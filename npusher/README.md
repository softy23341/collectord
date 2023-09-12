# README
NPusherD - сервис для отправки push-уведомлений.

Текущая версия поддерживает отправку через AWS SNS. И принимает задачи на отправку через RabbitMQ.

## Сборка
```bash
# server
make build

# test client
make build-mq-clietn

ls bin/
npusherd  npusher-mq-client
```

npusherd -- сервер.

npusher-mq-client -- клиент для добавления задач на push в mq-очередь.

## Запуск
Сервер:
```bash
bin/npusherd --config /tmp/cfg.toml --stdoutloglvl "debug"
```

Пример конфига можно посмотреть в корне репозитория (config-example.toml)

Клиент:
```bash
bin/npusher-mq-client --token "e4a0fac4 d2d0348e 830287b2 f5312f0e 22d001a3 0863e649 2f77359d 9a9e42fc" --message "test push" --sandbox
```
