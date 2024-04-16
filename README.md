# Планировщик задач
## Golang project
Планировщик хранит задачи, каждая из них содержит дату дедлайна и заголовок с комментарием. Задачи могут повторяться по заданному правилу: например, ежегодно, через какое-то количество дней, в определённые дни месяца или недели. Если отметить такую задачу как выполненную, она переносится на следующую дату в соответствии с правилом. Обычные задачи при выполнении будут просто удаляться.   
API содержит следующие операции:
* добавить задачу;
* получить список задач;
* удалить задачу;
* получить параметры задачи;
* изменить параметры задачи;
* отметить задачу как выполненную.