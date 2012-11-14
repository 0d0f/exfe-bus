import asyncore, socket

class Client(asyncore.dispatcher_with_send):
    def __init__(self, host, port, message):
        asyncore.dispatcher.__init__(self)
        self.create_socket(socket.AF_INET, socket.SOCK_STREAM)
        self.connect((host, port))
        self.out_buffer = message

    def handle_close(self):
        self.close()

    def handle_read(self):
        print 'Received', self.recv(1024)
        self.close()

c = Client('192.168.1.100', 8080, 'Hello, world')
asyncore.loop()
