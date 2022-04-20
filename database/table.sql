DROP IF EXISTS chat;
CREATE database chat;
use chat;

CREATE TABLE user (
    id int primary key auto_increment,
    username varchar(33) not null,
    password varchar(60) not null
);

CREATE TABLE room (
    id int primary key auto_increment,
    room_number varchar(30)
);

CREATE TABLE message (
  id int primary key auto_increment,
  owner_id int not null,
  room_id int not null,
  created_at datetime not null,
  content varchar(300),
  foreign key (owner_id) references user(id),
  foreign key (room_id) references room(id)
);