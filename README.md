### Video merger with ffmpeg

Запуск проекта:
```
docker-comose up -d
```

Приложение работает на порту 4000. Принимает 2 файла формата mp4.

Эндпоинт загрузки принимает следующие параметры (query params):

`localhost:4000/api/v1/video/upload`
- id - id операции;
- user_id - id юзера;
- video_id - id видео (1 или 2, соответственно первое и второе видео);

Эндпоинт выгрузки видео принимает следующие параметры (query params):

`localhost:4000/api/v1/video/download`
- id - id операции;
- user_id - id юзера;
- mod - желаемый вариант смерживания (pip или sbs). pip - картинка в картинке, sps - сторона к стороне;
