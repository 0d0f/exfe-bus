on run argv
	tell application "Messages"
		set recipient to item 1 of argv
		set content to item 2 of argv
		set sender to "E:_imessage@exfe.com"
		
		try
		send content to buddy recipient of service sender
		on error
		show chat chooser for buddy recipient of service sender
		send content to buddy recipient of service sender
		end
	end tell
end run