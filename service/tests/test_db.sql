BEGIN ;
INSERT INTO User (uuid, username) VALUES
	('69b618ef-2378-463a-b44e-f834513dff08', 'Mario'),
	('ede05cf3-12b5-491f-84e1-b925a52f6fde', 'Luigi'),
	('67a37587-4514-499f-b48e-e1bc49897a3b', 'ToadGiallo'),
	('350069d2-5e75-49c8-b6ef-d3e09ba131bc', 'ToadBlu'),
	('eb597bf9-7836-448e-8081-f546715c0955', 'Pitch');

INSERT INTO Conversation VALUES (1), (2), (3), (4), (5);

INSERT INTO Chat VALUES
    (1, '69b618ef-2378-463a-b44e-f834513dff08', 'ede05cf3-12b5-491f-84e1-b925a52f6fde'),
    (2, '67a37587-4514-499f-b48e-e1bc49897a3b', '69b618ef-2378-463a-b44e-f834513dff08'),
    (3, 'eb597bf9-7836-448e-8081-f546715c0955', '69b618ef-2378-463a-b44e-f834513dff08'),
    (4, 'ede05cf3-12b5-491f-84e1-b925a52f6fde', '350069d2-5e75-49c8-b6ef-d3e09ba131bc');

INSERT INTO GroupConversation (id, name, author) VALUES
    (5, 'Festa per pitch', '69b618ef-2378-463a-b44e-f834513dff08');

INSERT INTO User_Conversation VALUES
    ('67a37587-4514-499f-b48e-e1bc49897a3b', 5),
    ('350069d2-5e75-49c8-b6ef-d3e09ba131bc', 5),
    ('ede05cf3-12b5-491f-84e1-b925a52f6fde', 5),
    ('eb597bf9-7836-448e-8081-f546715c0955', 5);

INSERT INTO Message (id, conversation, author, sendAt, replyTo, content, attachment) VALUES
	(1, 1, '69b618ef-2378-463a-b44e-f834513dff08', '2010-05-28T15:36:56.200', NULL, 'Ciao Luigi, ho salvato la principessa, ti ho aggiunto ad un gruppo per la festa del salvataggio, quando puoi rispondi :D', NULL),
	(2, 1, 'ede05cf3-12b5-491f-84e1-b925a52f6fde', '2010-05-28T15:40:00.000', 1, 'Ho visto! Ti rispondo li!', NULL),
	(3, 3, 'eb597bf9-7836-448e-8081-f546715c0955', '2010-05-28T18:36:56.200', NULL, 'Grazie Mario, sono finalmente a casa... Mi mancava <3', NULL),
	(4, 5, '69b618ef-2378-463a-b44e-f834513dff08', '2010-05-28T13:36:56.200', NULL, 'Ciao, siete tutti invitati alla festa per Pitch domani alle 14.00 al livello 2-01!', NULL);

UPDATE User_Message SET status = 3 WHERE user = 'ede05cf3-12b5-491f-84e1-b925a52f6fde' AND message = 1;
UPDATE User_Message SET status = 2 WHERE user = '69b618ef-2378-463a-b44e-f834513dff08' AND message = 2;
UPDATE User_Message SET status = 3 WHERE user = 'ede05cf3-12b5-491f-84e1-b925a52f6fde' AND message = 4;
COMMIT ;
