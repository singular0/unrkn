# unRKN #

> A furore Roskomnadzoris libera nos, Domine

Скрипт позволяет создать немного оптимизированный список IPv4 адресов для маршрутизации через VPN
на вашем роутере на основе выгрузок РКН.

Источник данных - [API antizapret.info](https://antizapret.info/api.php). Помните, что по умолчанию
они ограничивают вас 6000 запросами в сутки.

Скрипт по умолчанию выгружает все IP-only записи. Кроме того, можно указать список доменных имён,
IP адреса которых из реестра он так же включит в список. При сравнении с записями реестра, имена
будут превращены в маски вида `^(.*\.)?domain.tld$`.

Список адресов будет слегка оптимизирован: если адрес или подсеть входят или равны другому адресу
или подсети из результирующего списка, они не будут включены повторно.

Список может быть выгружен в двух форматах:

* Один адрес/подсеть на строку (`raw`).
* Список команд вида `/ip firewall address-list add ...` для Mikrotik (`routeros`). Удобно заливать
  в конфигурацию командой `/import file-name=...`.

Скрипт написан на Go, собирается при помощи make.

## Командная строка ##

`unrkn [-a list_name] [-f format] [-o file] [-w file]`

| Ключ           | Описание |
|----------------|----------|
| `-a list_name` | При выводе списка в формате `routeros` добавить имя addres list'а `list={list_name}`. |
| `-f format`    | Формат выгрузки (см. выше): `raw` или `routeros`. По умолчанию `raw`. |
| `-o file`      | Записывать в указанный файл. Если опция не используется - выводить в stdout. |
| `-w file`      | Имя файла со списком доменных имён (см. выше). |

## Благодарности ##

Большой привет и спасибо [Филу](https://usher2.club) за его работу.
