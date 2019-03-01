TRUNCATE app CASCADE;
TRUNCATE userapp CASCADE;

INSERT INTO app(link, name, name_de, name_ru, image, about, about_de, about_ru, category)
    VALUES (
        '/terminal',
        'Terminal',
        'Die Konsole',
        'Командная строка',
        'img/app_covers/terminal.jpg',
        'Simple console, just like in your favorite Ubuntu',
        'Einfache Konsole, genau wie in Ihrem Lieblings Ubuntu',
        'Простая консоль, прямо как в любимой Убунте',
        '2018_2'
    );

INSERT INTO app(link, name, name_de, name_ru, image, about, about_de, about_ru, category)
    VALUES (
        '/snake',
        'Snake',
        'Snake',
        'Snake',
        'img/app_covers/snake.jpg',
        'Simple Snake game',
        'Ein Spiel ohne Schnickschnack, genau auf dem Handy',
        'Прямо как на Nokia 3310',
        '2018_2'
    );

INSERT INTO app(link, url, name, name_de, name_ru, image, about, about_de, about_ru, category)
    VALUES (
        '/proto',
        'https://rasseki.pro/',
        'Quizzy',
        'Quizzy',
        'Quizzy',
        'img/app_covers/proto.jpg',
        'Simple duel quiz by __proto__',
        'Duel quiz von __proto__',
        'Дуэль-викторина от __proto__',
        '2018_2'
    );

INSERT INTO app(link, url, name, name_de, name_ru, image, about, about_de, about_ru, category)
    VALUES (
        '/blep',
        'https://blep.me',
        'Blep',
        'Blep',
        'Blep',
        'img/app_covers/blep.jpg',
        'Elegant game about imagination by Stacktivity',
        'Elegantes Spiel um die Fantasie von Stacktivity',
        'Элегантная игра про воображение от Stacktivity',
        '2018_2'
    );

INSERT INTO app(link, url, name, name_de, name_ru, image, about, about_de, about_ru, category)
    VALUES (
        '/catchthealien',
        'https://itberries-frontend.herokuapp.com/',
        'Catch the alien!',
        'Fange den Alien!',
        'Поймай пришельца!',
        'img/app_covers/catch.jpg',
        'Dont let the alien leave the field! Game by ItBerries',
        'Lass den Alien nicht das Feld verlassen!',
        'Не позвольте пришельцу покинуть поле!',
        '2018_1'
    );

INSERT INTO app(link, url, name, name_de, name_ru, image, about, about_de, about_ru, category)
    VALUES (
        '/guardians',
        'https://chunk-frontend.herokuapp.com/',
        'Guardians',
        'Die Wächter',
        'Стражники',
        'img/app_covers/guardians.jpg',
        'Unusual 3D multi-player puzzle. Game by Chunk',
        'Ungewöhnliches 3D Multiplayer Puzzle. Spiel von Chunk',
        'Необычная 3D головоломка для нескольких игроков. Игра от команды Chunk',
        '2017_2'
    );

INSERT INTO app(link, url, name, name_de, name_ru, image, about, about_de, about_ru, category)
    VALUES (
        '/rhytmblast',
        'https://glitchless.surge.sh/',
        'Rhythm Blast',
        'Rhythm Blast',
        'Rhythm Blast',
        'img/app_covers/rhytmblast.jpg',
        'Space arcanoid on steroids. Game by Glitchless',
        'Arcade-Spiel im Weltraum von Glitchless',
        'Космический арканоид от Glitchless',
        '2017_2'
    );

INSERT INTO app(link, url, name, name_de, name_ru, image, about, about_de, about_ru, category)
    VALUES (
        '/ketnipz',
        'https://playketnipz.ru/',
        'Ketnipz',
        'Ketnipz',
        'Ketnipz',
        'img/app_covers/ketnipz.jpg',
        'Arcade game with cartoony graphics by DeadMolesStudio',
        'Arcade-Spiel mit schöner Grafik von DeadMolesStudio',
        'Аркадная игра с приятной графикой от DeadMolesStudio',
        '2018_2'
    );

INSERT INTO app(link, url, name, name_de, name_ru, image, about, about_de, about_ru, category)
    VALUES (
        '/kekmate',
        'https://kekmate.tech/',
        'Chessmate',
        'Chessmate',
        'Шахматы',
        'img/app_covers/chess.jpg',
        'Classic chess with advanced AI game by Parashutnaya Molitva',
        'Klassisches Schachspiel mit fortgeschrittener KI von Paraschyutnaya Molitva',
        'Классические шахматы с продвинутым ИИ от Парашютной Молитвы',
        '2018_2'
    );

INSERT INTO app(link, url, name, name_de, name_ru, image, about, about_de, about_ru, category)
    VALUES (
        '/rpsarena',
        'http://rpsarena.ru',
        'RPS Arena',
        'Schere, Stein, Papier Arena',
        'RPS Arena',
        'img/app_covers/rps.jpg',
        'Multiplayer version of classic Rock–Paper–Scissors game by 42',
        'Schere, Stein, Papier für zwei Spieler des Teams 42',
        'Камень-ножницы-бумага для двух игроков от команды 42',
        '2018_2'
    );


-- don't work
-- INSERT INTO app(link, url, name, image, about, category) VALUES (
        -- '/simplegame',
-- '    https://simplegame.ru.com/',
-- '    Simple Game',
-- '    img/app_covers/simplegame.jpg',
-- '    Game by Simple Name',
-- '    2018_2' );
-- -- INSERT INTO app(link, url, name, image, about, category) VALUES (
        -- '/yetanothergame',
-- '    https://yet-another-game.ml/',
-- '    Yet Another Game',
-- '    img/app_covers/yetanothergame.jpg',
-- '    Game by Yet Another Game',
-- '    2018_2' );
-- -- INSERT INTO app(link, url, name, image, about, category) VALUES (
        -- '/tron',
-- '    https://codeloft.ru',
-- '    Tron: Remastered',
-- '    img/app_covers/tron.jpg',
-- '    Game by codeloft',
-- '    2018_2' );e_ru, image, about, about_de, about_ru, category)
