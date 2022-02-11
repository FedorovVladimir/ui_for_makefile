ifeq (command_with_arg,$(firstword $(MAKECMDGOALS)))
  arg := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  $(eval $(arg):;@:)
endif

command_1:
	echo "command_1" > text.txt

# command 1 description
command_2:
	echo "command_2" > text.txt

command_3:
	echo "command_3" > text.txt

.PHONY: command_with_arg
command_with_arg:
	echo "command_with_arg:" $(arg) > text.txt
