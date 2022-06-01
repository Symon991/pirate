Simple console application to search for torrents using the api behind the pirate bay site.
--It can build magnet links and add them directly to a qBittorrent instance.

Example usage:

There are 3 subcommands

torrent: to search torrents and add them to a qBittorrent instance
--subtitle: to search subtitles
--config: to edit the persistent configuration

pirate.exe config -url "192.168.1.10:8080" -name "home"
--pirate.exe torrent -s "Ambulance 2022 2160p" -add "home" -c "Film"
--pirate.exe subtitle -s "the batman" -l eng


