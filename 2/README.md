# mycut

Распределённая CLI-утилита `cut` на Go с quorum и concurrency.

## Возможности

- локальный режим, аналогичный `cut`;
- распределённый режим через worker/client;
- обработка shard-файлов;
- quorum на клиенте;
- concurrency внутри worker через goroutines и channels.

## Флаги

- `-mode` — режим: `local`, `worker`, `client`
- `-f` — какие поля вывести
- `-d` — разделитель полей
- `-s` — выводить только строки, содержащие разделитель
- `-addr` — адрес worker
- `-file` — shard-файл worker
- `-servers` — список серверов client через запятую
- `-quorum` — кворум для client

## Сборка

```powershell
go mod init mycut
go build -o mycut.exe
```
```Тестовые shard-файлы
shard1.txt
shard2.txt
shard3.txt
```

## Примеры использования
Local mode
```
Get-Content .\shard1.txt | .\mycut.exe -mode=local -f "1,3" -d ","
```
```
Ожидаемый вывод
john,london
alice,paris
```
Worker mode
```
.\mycut.exe -mode=worker -addr="localhost:9001" -file="shard1.txt"
.\mycut.exe -mode=worker -addr="localhost:9002" -file="shard1.txt"
.\mycut.exe -mode=worker -addr="localhost:9003" -file="shard1.txt"
```
Client mode
```
.\mycut.exe -mode=client -servers="localhost:9001,localhost:9002,localhost:9003" -quorum=2 -f "1,3" -d ","
```
```
Ожидаемый вывод
john,london
alice,paris
bob,berlin
kate,madrid
tom,rome
eva,vienna
```
Проверка quorum
Доступны 2 worker из 3
```
.\mycut.exe -mode=client -servers="localhost:9001,localhost:9002,localhost:9003" -quorum=2 -f "1,3" -d ","
```
```
Ожидаемый вывод 
ошибка насчет 9001
john,london
alice,paris
bob,berlin
kate,madrid
```
Доступен 1 worker из 3
```
.\mycut.exe -mode=client -servers="localhost:9001,localhost:9002,localhost:9003" -quorum=2 -f "1,3" -d ","
```
```Ожидаемый вывод
ошибка насчет 9001,9002
quorum not reached: 1 / 2
```
Локальный режим

Оригинальная cut в Git Bash / WSL:

```cat shard1.txt | cut -d "," -f 1,3```

Ожидаемый вывод:

```
john,london
alice,paris
```
mycut:

```
Get-Content .\shard1.txt | .\mycut.exe -mode=local -f "1,3" -d ","
```

Ожидаемый вывод:

```
john,london
alice,paris
```
Distributed mode

Для полного сравнения нужно использовать все shard и полный кворум.

Оригинальная cut для общего файла all.txt:

```
cat all.txt | cut -d "," -f 1,3
```

Ожидаемый вывод:

```
john,london
alice,paris
bob,berlin
kate,madrid
tom,rome
eva,vienna
```

mycut:

```
.\mycut.exe -mode=client -servers="localhost:9001,localhost:9002,localhost:9003" -quorum=3 -f "1,3" -d ","
```
Ожидаемый вывод:

```
john,london
alice,paris
bob,berlin
kate,madrid
tom,rome
eva,vienna
```