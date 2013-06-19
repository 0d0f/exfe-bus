#include "base.h"
#include "log.h"
#include "buffer.h"

#ifdef HAVE_MYSQL
#include <mysql.h>
#endif

#include "plugin.h"

#include <ctype.h>
#include <stdlib.h>
#include <string.h>

#define DEBUG(...)                                           \
        log_error_write(srv, __FILE__, __LINE__, __VA_ARGS__);

#define HEADER(con, key)                                                \
    (data_string *)array_get_element((con)->request.headers, (key))

/**
 * this is a mysql_obj for a lighttpd plugin
 *
 * just replaces every occurance of 'mysql_obj' by your plugin name
 *
 * e.g. in vim:
 *
 *   :%s/mysql_obj/myhandler/
 *
 */



/* plugin config for all request/connections */

typedef struct {
	array *query;
	buffer *table;
	buffer *key;
	array *map;
	array *extra;
	buffer *sql_query;
} plugin_config;

typedef struct {
	PLUGIN_DATA;

	plugin_config **config_storage;

	MYSQL *mysql;
	plugin_config conf;
} plugin_data;

/* init the plugin data */
INIT_FUNC(mod_mysql_obj_init) {
	plugin_data *p;
	p = calloc(1, sizeof(*p));
	return p;
}

/* detroy the plugin data */
FREE_FUNC(mod_mysql_obj_free) {
	plugin_data *p = p_d;

	UNUSED(srv);

	if (!p) return HANDLER_GO_ON;

	if (p->config_storage) {
		size_t i;

		for (i = 0; i < srv->config_context->used; i++) {
			plugin_config *s = p->config_storage[i];

			if (!s) continue;

			array_free(s->query);
			buffer_free(s->table);
			buffer_free(s->key);
			array_free(s->map);
			array_free(s->extra);
			buffer_free(s->sql_query);

			free(s);
		}
		free(p->config_storage);
	}

	mysql_close(p->mysql);
	free(p);

	return HANDLER_GO_ON;
}

/* handle plugin config and check values */

SETDEFAULTS_FUNC(mod_mysql_obj_set_defaults) {
	plugin_data *p = p_d;
	size_t i = 0;

	config_values_t cv[] = {
		{ "mysql_obj.host",             NULL, T_CONFIG_STRING, T_CONFIG_SCOPE_SERVER },       /* 0 */
		{ "mysql_obj.port",             NULL, T_CONFIG_SHORT,  T_CONFIG_SCOPE_SERVER },       /* 1 */
		{ "mysql_obj.user",             NULL, T_CONFIG_STRING, T_CONFIG_SCOPE_SERVER },       /* 2 */
		{ "mysql_obj.pass",             NULL, T_CONFIG_STRING, T_CONFIG_SCOPE_SERVER },       /* 3 */
		{ "mysql_obj.sock",             NULL, T_CONFIG_STRING, T_CONFIG_SCOPE_SERVER },       /* 4 */
		{ "mysql_obj.db",               NULL, T_CONFIG_STRING, T_CONFIG_SCOPE_SERVER },       /* 5 */
		{ "mysql_obj.insert",           NULL, T_CONFIG_LOCAL,  T_CONFIG_SCOPE_CONNECTION },   /* 6 */
		{ NULL,                         NULL, T_CONFIG_UNSET,  T_CONFIG_SCOPE_UNSET }
	};

	if (!p) return HANDLER_ERROR;

	buffer *host = buffer_init();
	unsigned short port = 0;
	buffer *user = buffer_init();
	buffer *pass = buffer_init();
	buffer *sock = buffer_init();
	buffer *db = buffer_init();
	p->config_storage = calloc(srv->config_context->used, sizeof(specific_config *));

	for (i = 0; i < srv->config_context->used; i++) {
		array *insert = array_init();	
		cv[0].destination = host;
		cv[1].destination = &port;
		cv[2].destination = user;
		cv[3].destination = pass;
		cv[4].destination = sock;
		cv[5].destination = db;
		cv[6].destination = insert;

		if (0 != config_insert_values_global(srv, ((data_config *)srv->config_context->data[i])->value, cv)) {
			return HANDLER_ERROR;
		}

		array *ca = ((data_config *)srv->config_context->data[i])->value;
		data_array *du = (data_array *)array_get_element(ca, "mysql_obj.insert");

		if (du == NULL) continue;
		if (du->type != TYPE_ARRAY) return HANDLER_ERROR;

		data_array *query = (data_array *)array_get_element(du->value, "query");
		data_string *table = (data_string *)array_get_element(du->value, "table");
		data_string *key = (data_string *)array_get_element(du->value, "key");
		data_array *map = (data_array *)array_get_element(du->value, "map");
		data_array *extra = (data_array *)array_get_element(du->value, "extra");

		plugin_config *s;
		s = calloc(1, sizeof(*s));
		p->config_storage[i] = s;

		s->query = array_init();
		s->table = buffer_init();
		s->key = buffer_init();
		s->map = array_init();
		s->extra = array_init();
		s->sql_query = buffer_init_string("SELECT ");

		for (uint j=0; j<query->value->used; j++) {
			data_unset *token = query->value->data[j]->copy(query->value->data[j]);
			array_insert_unique(s->query, token); 
		}
		buffer_copy_string(s->table, table->value->ptr);
		buffer_copy_string(s->key, key->value->ptr);
		for (uint j=0; j<map->value->used; j++) {
			buffer *mkey = map->value->data[j]->key;
			data_string *mvalue = (data_string *)map->value->data[j];
			array_set_key_value(s->map, mkey->ptr, strlen(mkey->ptr), mvalue->value->ptr, strlen(mvalue->value->ptr)); 
			buffer_append_string(s->sql_query, "`");
			buffer_append_string(s->sql_query, mkey->ptr);
			buffer_append_string(s->sql_query, "`");
			if (j != map->value->used - 1) {
				buffer_append_string(s->sql_query, ", ");
			}
		}
		buffer_append_string(s->sql_query, " FROM `");
		buffer_append_string_buffer(s->sql_query, s->table);
		buffer_append_string(s->sql_query, "` WHERE `");
		buffer_append_string_buffer(s->sql_query, s->key);
		buffer_append_string(s->sql_query, "`=\"");
		
		for (uint j=0; j<extra->value->used; j++) {
			buffer *mkey = extra->value->data[j]->key;
			data_string *mvalue = (data_string *)extra->value->data[j];
			array_set_key_value(s->extra, mkey->ptr, strlen(mkey->ptr), mvalue->value->ptr, strlen(mvalue->value->ptr)); 
		}
		array_free(insert);
	}

	if (NULL == (p->mysql = mysql_init(NULL))) {
		log_error_write(srv, __FILE__, __LINE__, "s", "mysql_init() failed, exiting...");
		return HANDLER_ERROR;
	}
#if MYSQL_VERSION_ID >= 50013
	/* in mysql versions above 5.0.3 the reconnect flag is off by default */
	my_bool reconnect = 1;
	mysql_options(p->mysql, MYSQL_OPT_RECONNECT, &reconnect);
#endif

#if MYSQL_VERSION_ID >= 40100
	/* CLIENT_MULTI_STATEMENTS first appeared in 4.1 */ 
	if (!mysql_real_connect(p->mysql, host->ptr, user->ptr, pass->ptr,
		db->ptr, port, sock->ptr, CLIENT_MULTI_STATEMENTS)) {
#else
	if (!mysql_real_connect(p->mysql, host->ptr, user->ptr, pass->ptr,
		db->ptr, port, sock->ptr, 0)) {
#endif
		log_error_write(srv, __FILE__, __LINE__, "ss", "mysql connect error:", mysql_error(p->mysql));

		return HANDLER_ERROR;
	}

#ifdef FD_CLOEXEC
	fcntl(p->mysql->net.fd, F_SETFD, FD_CLOEXEC);
#endif

	buffer_free(host);
	buffer_free(user);
	buffer_free(pass);
	buffer_free(sock);
	buffer_free(db);

	return HANDLER_GO_ON;
}

#define PATCH(x) \
	p->conf.x = s->x;
static int mod_mysql_obj_patch_connection(server *srv, connection *con, plugin_data *p) {
	size_t i, j;
	plugin_config *s = p->config_storage[0];

	if (s != NULL) {
		PATCH(query);
		PATCH(table);
		PATCH(key);
		PATCH(map);
		PATCH(extra);
		PATCH(sql_query);
	} else {
		p->conf.query = NULL;
		p->conf.table = NULL;
		p->conf.key = NULL;
		p->conf.map = NULL;
		p->conf.extra = NULL;
		p->conf.sql_query = NULL;
	}

	/* skip the first, the global context */
	for (i = 1; i < srv->config_context->used; i++) {
		data_config *dc = (data_config *)srv->config_context->data[i];
		s = p->config_storage[i];

		/* condition didn't match */
		if (!config_check_cond(srv, con, dc)) continue;

		/* merge config */
		for (j = 0; j < dc->value->used; j++) {
			data_unset *du = dc->value->data[j];

			if (buffer_is_equal_string(du->key, CONST_STR_LEN("mysql_obj.insert"))) {
				PATCH(query);
				PATCH(table);
				PATCH(key);
				PATCH(map);
				PATCH(extra);
				PATCH(sql_query);
			}
		}
	}

	return 0;
}
#undef PATCH

URIHANDLER_FUNC(mod_mysql_obj_uri_handler) {
	plugin_data *p = p_d;

	UNUSED(srv);
	if (con->mode != DIRECT) return HANDLER_GO_ON;
	if (con->uri.query->used == 0) return HANDLER_GO_ON;

	mod_mysql_obj_patch_connection(srv, con, p);

	if (p->conf.query == NULL) {
		return HANDLER_GO_ON;
	}

	char *token = NULL;
	int token_len = -1;
	for (uint i = 0; i < p->conf.query->used; i++) {
		buffer *key = buffer_init_string(((data_string *)p->conf.query->data[i])->value->ptr);
		buffer_append_string(key, "=");
		token = strstr(con->uri.query->ptr, key->ptr);
		int key_len = strlen(key->ptr);
		buffer_free(key);
		if (token != NULL) {
			token += key_len;
			char *end = strchr(token, '&');
			if (end != NULL) {
				token_len = end - token;
			} else {
				token_len = con->uri.query->ptr + strlen(con->uri.query->ptr) - token;
			}
			break;
		}
	}
	if (token == NULL) {
		return HANDLER_GO_ON;
	}
	buffer *query = buffer_init_string(p->conf.sql_query->ptr);
	buffer_append_string_len(query, token, token_len);
	buffer_append_string(query, "\"");
        if (mysql_query(p->mysql, query->ptr)) {
                log_error_write(srv, __FILE__, __LINE__, "ssss", "query", query->ptr, "error:", mysql_error(p->mysql));
		buffer_free(query);
		return HANDLER_GO_ON;
        }
	buffer_free(query);

        MYSQL_RES *result = mysql_store_result(p->mysql);
	unsigned cols = mysql_num_fields(result);
	MYSQL_ROW row = mysql_fetch_row(result);
	unsigned long *lengths = mysql_fetch_lengths(result);

	if (!row || cols != p->conf.map->used ) {
		mysql_free_result(result);
#if MYSQL_VERSION_ID >= 40100
		while (mysql_next_result(p->mysql) == 0);
#endif
		return HANDLER_GO_ON;
	}

	for (uint i = 0; i < p->conf.map->used; i++) {
		data_string *key = (data_string *)p->conf.map->data[i];
		array_set_key_value(con->request.headers, key->value->ptr, key->value->used-1, row[i], lengths[i]);
	}

	for (uint i = 0; i < p->conf.extra->used; i++) {
		buffer *key = p->conf.extra->data[i]->key;
		data_string *value = (data_string *)p->conf.extra->data[i];
		array_set_key_value(con->request.headers, key->ptr, key->used-1, value->value->ptr, value->value->used-1);
	}

	mysql_free_result(result);
#if MYSQL_VERSION_ID >= 40100
	while (mysql_next_result(p->mysql) == 0);
#endif

	for (uint i = 0; i<con->request.headers->used; i++) {
		data_unset *d = con->request.headers->data[i];
	}

	return HANDLER_GO_ON;
}

/* this function is called at dlopen() time and inits the callbacks */

int mod_mysql_obj_plugin_init(plugin *p) {
	p->version     = LIGHTTPD_VERSION_ID;
	p->name        = buffer_init_string("mysql_obj");

	p->init        = mod_mysql_obj_init;
	p->handle_uri_clean  = mod_mysql_obj_uri_handler;
	p->set_defaults  = mod_mysql_obj_set_defaults;
	p->cleanup     = mod_mysql_obj_free;

	p->data        = NULL;

	return 0;
}
