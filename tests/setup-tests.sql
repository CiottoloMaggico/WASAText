BEGIN TRANSACTION;

INSERT INTO User (uuid, username)
VALUES ('TWFyaW8=', 'Mario'),
	   ('VG9hZEdpYWxsbw==', 'ToadGiallo'),
	   ('VG9hZEJsdQ==', 'ToadBlu'),
	   ('THVpZ2k=', 'Luigi');

INSERT INTO Conversation
VALUES (1),
	   (2),
	   (3);

INSERT INTO GroupConversation (id, name)
VALUES (3, 'Festa di Mario');

INSERT INTO Chat
VALUES (1, 'TWFyaW8=', 'THVpZ2k='),
	   (2, 'VG9hZEdpYWxsbw==', 'THVpZ2k=');

INSERT INTO User_Conversation
VALUES ('TWFyaW8=', 1),
	   ('THVpZ2k=', 1),
	   ('VG9hZEdpYWxsbw==', 2),
	   ('THVpZ2k=', 2),
	   ('TWFyaW8=', 3),
	   ('VG9hZEdpYWxsbw==', 3),
	   ('VG9hZEJsdQ==', 3),
	   ('THVpZ2k=', 3);
;

INSERT INTO Message
VALUES (1, 1, 'THVpZ2k=',
		'2008-11-11 13:23:44', NULL,
		'Sono contento che tu abbia salvato la principessa',
		NULL),
	   (2, 1, 'TWFyaW8=',
		'2008-11-11 13:23:45', NULL,
		'Grazie Luigi, domani faremo una festa vuoi venire? Ci trovi al livello 1-2 alle 17',
		NULL),
	   (3, 3, 'TWFyaW8=', '2008-11-11 13:26:43', NULL,
		'Siete tutti invitati alla festa per Pitch al livello 1-2 alle 17',
		NULL),
	   (4, 2, 'VG9hZEdpYWxsbw==', '2008-11-11 13:23:43', NULL, 'Vieni alla festa di mario?', NULL);


INSERT INTO User_Message
VALUES (1, 'TWFyaW8=', 'seen', NULL),
	   (1, 'THVpZ2k=', 'seen', NULL),
	   (2, 'TWFyaW8=', 'seen', NULL),
	   (2, 'THVpZ2k=', 'delivered', NULL),
	   (4, 'VG9hZEdpYWxsbw==', 'seen', NULL),
	   (4, 'THVpZ2k=', 'sent', NULL),
	   (3, 'VG9hZEJsdQ==', 'sent', NULL),
	   (3, 'TWFyaW8=', 'seen', NULL),
	   (3, 'THVpZ2k=', 'delivered', NULL),
	   (3, 'VG9hZEdpYWxsbw==', 'sent', NULL);
END TRANSACTION;
