# Xsolla
Тестовое задание для Xsolla, основа - это файл "Test task for Middle v2".

В рамках задания необходимо реализовать сервис доставки пиццы.
Предполагаемая реализация представляет 3 отдельных микросервиса - Shop, Kitchen, Delivery.
Каждый из них включает в себе собственную бизнес-логику, с доменным структурами, взаимодействием с клиентом, репозиторий. 
Межсервисное взимодействие осуществляется асинхронно, через брокера сообщений NATS.

Общий процесс 

Сервис Shop:
Создание заказа клиентом (MakeOrder, статус New) + задача -> 
Вывод созданных заказов для подтверждения (ListOrders) ->
Подтверждение заказа администратором (ChangeOrderStatus, статус Confirmed) + задача ->
По факту получения события об обновлении Cooking на InProgress обновление заказа на Cooking + задача ->
По факту получения события об обновлении Cooking на Completed обновление заказа на Cooked + задача ->
По факту получения события об обновлении Delivery на InProgress обновление заказа на Delivering + задача ->
По факту получения события об обновлении Delivery на Completed обновление заказа на Completed + задача ->

Сервис Kitchen 
По факту получения события о создании заказа происходит создание Cooking, статус New ->
По факту получения события об обновлении статуса заказа до Confirmed происходит обновлении статуса
Cooking до NeedToStart ->
Вывод Cooking с NeedToStart для поваров (ListOfCooking) -> 
По факту принятия заказа поваром обновление статуса Cooking на InProgress (ChangeCookingStatus) + задача об обновлении в очередь в транзакции ->
По факту завершения заказа поваром обновление статуса Cooking на Completed (ChangeCookingStatus) + задача об обновлении в очередь в транзакции

Сервис Delivery
По факту получения события о создании заказа происходит создание Delivery, статус New ->
По факту получения события об обновлении статуса заказа до Cooked происходит обновлении статуса
Delivery до NeedToGo ->
Вывод Cooking с NeedToGo для курьеров (ListOfDelivery) -> 
По факту начала процесса доставки курьером происходит обновление статуса Delivery на InProgress(ChangeDeliveryStatus) + задача об обновлении в очередь в транзакции ->
По факту завершения процесса доставки курьером происходит обновление статуса Delivery на Completed (ChangeDeliveryStatus) + задача об обновлении в очередь в транзакции

Комментарии относительно пунктов к описанию:

todo