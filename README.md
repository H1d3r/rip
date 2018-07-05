# RIP
Простой DNS сервер с возможность кодирования IP адреса в имени.

# Кодирование IP
DNS сервер парсит имя, пытается вычленить из него требуемый IP и формирует ответ.
Правила:
```
    <optional-prefix>.<IPv4>.4.<zone>  -> returns A record with <IPv4> address
    <optional-prefix>.<IPv6>.6.<zone>  -> returns AAAA record with <IPv6> address
    <proxy-name>.p.<zone>  -> resolve proxy name and returns it
    <ip>.<ip>.r.<zone>  -> pick random IP
    <ip1>.<ip2>.l.<zone>  -> loop over <ip1> and <ip2>
    <cname>.c.<zone>  -> return CNAME record with <cname>
    <any-name>.<zone>  -> returns default address
```

# Формат IP
IP может быть представлен в двух вариантах - dash-delimited и base-16.
К примеру, эти записи `0a000001` и `10-0-0-1` эквивалентны и указывают `10.0.0.1`

# Формат имен
cname и proxy поддерживают две логики резолвинга имен - префиксный и dash-delimited.
К примеру:
```
    something.victim.com.c.evil.com -> CNAME to something.victim.com
    something.victim-com.c.evil.com -> CNAME to victim.com
```

# Пример сервера
Запустим NS сервер дл зоны `example.com` с дефолтными IP `77.88.55.70` и `2a02:6b8:a::a`:
```
$ rip ns --zone=example.com --ipv4=77.88.55.70 --ipv6=2a02:6b8:a::a
```
При запросе к нему, мы должны получить следующие ответы:
```
# IPv4
    1-1-1-1.4.example.com ->  1.1.1.1  && 2a02:6b8:a::a
    foo.1-1-1-1.4.example.com ->  1.1.1.1  && 2a02:6b8:a::a
    bar.foo.1-1-1-1.4.example.com ->  1.1.1.1  && 2a02:6b8:a::a

# IPv6
    2a01-7e01--f03c-91ff-fe3b-c9ba.6.example.com    ->  2a01:7e01::f03c:91ff:fe3b:c9ba  && 77.88.55.70
    foo.2a01-7e01--f03c-91ff-fe3b-c9ba.6.example.com    -> 2a01:7e01::f03c:91ff:fe3b:c9ba  && 77.88.55.70
    foo.--1.6.example.com   ->  ::1

# Random
    0a000002.0a000001.r.example.com ->  random between 10.0.0.1 and 10.0.0.2

# Loop
    8ba299a7.8ba299a8.l.example.com ->  loop over 139.162.153.167 and 139.162.153.168

# Cname
    ya.ru.c.example.com ->  canonical name ya.ru
    google.com.c.example.com ->  canonical name google.com

# Proxy
    ya.ru.p.example.com ->  87.250.250.242 && 2a02:6b8::2:242
    google.com.exampl.ecom  ->  64.233.164.102 and 2a00:1450:4010:c07::64
```

# IP-конвертер
Если вы такая же ленивая жопа как и я - `rip` умеет конвертировать IP-адреса в base-16 форму:
```
$ rip ip2hex fe80::fa94:c2ff:fee5:3cf6 127.0.0.1
fe80000000000000fa94c2fffee53cf6	7f000001
```