CREATE DATABASE manga_library;
CREATE TABLE manga_library.users
(
    id int(11) not null auto_increment,
    name varchar (255) not null,
    email varchar (255) default '',
    password varchar(512) not null ,
    api_key varchar(512) not null,
    current_jobs int(10) default 0,
    max_jobs int(10) default 0,
    age int(10) default 18,
    last_login timestamp default now(),
    is_active bool default 1,
    is_admin bool default 0,
    primary key (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=COMPRESSED;

INSERT INTO manga_library.users(name,email,password,api_key,is_admin)VALUES ("Seann Moser","seannsea@gmail.com","","12345",1);


CREATE TABLE manga_library.sites
(
    id int(11) not null auto_increment,
    name varchar(512) not null ,
    base_url varchar(512) not null ,
    search_url varchar(512) not null ,
    base_path blob,
    use_sub_path bool default false,
    meta_data mediumblob,
    min_age int default 0,
    primary key (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=COMPRESSED;


CREATE TABLE manga_library.web_component
(
    id int(11) not null auto_increment,
    site_id int(11) not null,
    name varchar(255) default '',
    tag varchar(128) default '',
    attribute varchar(255) default '',
    value varchar(255) default '',
    is_link bool default false,
    is_download bool default false,
    link_attributes blob,
    element_data longblob,
    parent int(11) default 0,
    reverse bool,
    delay int(10) default 5,
    meta_data blob,
    foreign key (site_id) REFERENCES manga_library.sites (id),
    primary key (id,site_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=COMPRESSED;


CREATE TABLE manga_library.jobs
(
    id int(11) not null auto_increment,
    user_id int(11) not null,
    site_id int(11) not null,
    name varchar(512) not null default '',
    job_context blob,
    start_time varchar(512) not null,
    #Progress
    current int(11) not null default 0,
    total int(11) not null default 0,
    est_finish varchar(512) not null,
    message blob,
    job_data longblob,
    primary key (id,user_id),
    foreign key (user_id) REFERENCES manga_library.users (id),
    foreign key (site_id) REFERENCES manga_library.sites (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=COMPRESSED;


CREATE TABLE manga_library.books
(
    id int(11) not null auto_increment,
    user_id int(11) not null,
    site_id int(11) not null,
    is_public bool default false,
    views int(32) default 0,
    downloads int(16) default 0,
    job_id int(11) not null,
    chapter varchar(1024) default '',
    volume varchar(1024) default '',
    name varchar(1024) default '',
    description blob,
    meta_data mediumblob,
    file_path varchar(512) default '',
    cover_img varchar(512) default '',
    pages int default 0,
    primary key (id,user_id),
    foreign key (user_id) REFERENCES manga_library.users (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=COMPRESSED;



CREATE TABLE manga_library.library
(
    user_id int(11) not null,
    book_id int(11) not null,
    collection varchar(512) default 'default',
    progress int(11) default 0,
    rating int(10) default 0,
    favorite bool default false,
    primary key (user_id,book_id),
    foreign key (user_id) REFERENCES manga_library.users (id),
    foreign key (book_id) REFERENCES manga_library.books (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=COMPRESSED;