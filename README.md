# WB Tech: level # 0 (Golang)

### Запуск проекта

### Способ 1: С использованием Docker (с Docker hub)

1. Скачайте контейнер Docker командой:

    ```
    docker pull happydayaway/wb_go_l0:latest
    ```

2. Запустите контейнер командой:

    ```
    docker run -d -p 80:1234 --rm --name wbapp happydayaway/wb_go_l0
    ```

### Способ 2: С использованием Docker (с Github)

1. Скачайте проект с GitHub командой:

    ```
    git clone https://github.com/karrrakotov/wb_go_lvl0.git
    ```

2. Откройте консоль, перейдите в корневую папку проекта и выполните команду:

    ```
    make build-compose
    ```
