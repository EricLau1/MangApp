create database dbmanga
default character set utf8
default collate utf8_general_ci;

create table if not exists mangas(
id int auto_increment primary key,
descricao varchar(255),
formato enum('tanko', 'meio-tanko', 'livro'),
quantidade int,
volumes varchar(255),
status char(1) default '0',
valor decimal(10,2) default 0.0
)
default charset=utf8;

insert into mangas (descricao, formato, quantidade, volumes, status) values
('Bakuman', 'tanko', 20, '1 ao 20', '1');