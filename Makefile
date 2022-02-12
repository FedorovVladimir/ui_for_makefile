## group: Простые

## команда с выбором
## спрашивает уверены ли вы в своем выборе
command_1:
	@./run.sh

## вывод в файл
command_2:
	echo "command_2" > text.txt

## group: С параметрами

## вывод текста в файл
command_with_arg_1:
	echo "command_with_arg:" $(arg) > text.txt

command_with_arg_2:
	echo "command_with_arg:" $(arg) > text.txt

command_with_arg_3:
	echo "command_with_arg:" $(arg) > text.txt
