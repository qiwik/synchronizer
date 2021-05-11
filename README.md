# Синхронизатор папок

Данное приложение выполняет фоновую синхронизацию двух папок, к которым имеется доступ с вашей
локальной машины. Приложение реализовано посредством сканирования файловых деревьев
указанных папок и их дальнейшего сравнения на наличие различий.

## Начальные параметры для работы

Абсолютные пути к необходимым папкам прописываются через командную строку, используя
следующий синтаксис:

`-source=... -copy=...`,

где первый параметр - это путь к папке, требующей синхронизации, а второй - путь к папке,
куда синхронизация будет происходить.

## Принцип работы

Приложение сравнивает в фоновом режиме наличие папок и файлов в двух файловых деревьях,
отслеживая изменения в них. Поддерживает следующие сценарии:

- копирование папок и файлов в директорию-копию при их полном отсутствии;
- сравнение хэш-сумм файлов и их замена в директории-копии при обнаружении изменений;
- удаление папок и файлов из директории-копии, если в синхронизируемой директории их
  не существует.

Выход из приложения осуществляется через комбинацию ctrl+c. При этом выход завершается
корректно.

Время обновления файловых деревьев приложение высчитывает само на основе данных о весе
обеих папок, чтобы дать возможность слабым машинам выполнить действия удаления/копирования
корректно. Вся информация о работе приложения записывается в лог файл, находящийся
в директории /cmd/app.

## Будущее изменение

На данный момент приложение завершает свою работу в случае, когда при сканировании
папок обнаруживается, что в одну из них происходит загрузка нового файла. Пока что данное
действие ломает выполнение приложения, поэтому прошу не загружать новые файлы в указанные
вами папки, если вы не уверенны в том, что машина справится с копированием до окончания
времени бездействия приложения