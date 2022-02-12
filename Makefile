# command_1 description 1
command_1:
	@./run.sh

# command_2 description 2
command_2:
	echo "command_2" > text.txt

# command_3 description 3
command_3:
	echo "command_3" > text.txt

# command_with_arg_1 description with args
# arg: input text
command_with_arg_1:
	echo "command_with_arg:" $(arg) > text.txt

# command_with_arg_2 description with args
# arg: arr=["1", "2", "3"]
command_with_arg_2:
	echo "command_with_arg:" $(arg) > text.txt

# command_with_arg_3 description with args
# arg: arrFrom=command_3
command_with_arg_3:
	echo "command_with_arg:" $(arg) > text.txt
