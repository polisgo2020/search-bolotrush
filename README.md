# search-bolotrush
# Текстовый поиск 
Данная программа составляет обратный индекс по словам, содержащимся в текстовых документах.
У пользователся есть возможность выбрать необходимое действие, введя нужные флаги при запуске программы. 
Следующим параметром при запуске является путь к текстовым файлам
### Запуск может производится со следующими флагами:
* `-f` - запись обратного индекса в отдельный текстовый файл
* `-s="пример запроса"` - вывод в командную стоку результата поиска самого подходящего документа по введенной поисковой фразе
* `-web="интерфейс"` - запуск сервера по адресу, который слушает указанный интерфейс и принимает http GET запрос, содержащий поисковую фразу и возвращает список подходящих документов 
