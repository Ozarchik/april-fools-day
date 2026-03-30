# April Fools Day

A small April Fools' Day site. Each visitor gets a unique sequential number. Returning visitors always see the same number.

## Stack

- Go (net/http)
- PostgreSQL
- Docker

## Running locally

```bash
docker compose up --build
```

Open [http://localhost:8081](http://localhost:8081).

## Deployment

```bash
git clone https://github.com/Ozarchik/april-fools-day.git
cd april-fools-day
docker compose up -d --build
```

Open [http://localhost:8081](http://localhost:8081).

## Apache reverse proxy

Enable required modules:

```bash
sudo a2enmod proxy proxy_http
sudo systemctl restart apache2
```

Add a virtual host config (e.g. `/etc/apache2/sites-available/april-fools.conf`):

```apache
<VirtualHost *:80>
    ServerName <your-domain>

    ProxyPass        / http://localhost:8081/
    ProxyPassReverse / http://localhost:8081/
</VirtualHost>
```

Enable and reload:

```bash
sudo a2ensite april-fools.conf
sudo systemctl reload apache2
```
