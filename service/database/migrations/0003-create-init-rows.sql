INSERT INTO MessageStatus (id, name)
VALUES (1, 'sent'),
	   (2, 'delivered'),
	   (3, 'seen')
ON CONFLICT (id) DO NOTHING;

INSERT INTO Image (uuid, extension, width, height, fullUrl)
VALUES
    ('default_group_image', '.jpg', 840, 880, '/media/images/default_group_image.jpg'),
	('default_user_image', '.jpg', 8000, 8000, '/media/images/default_user_image.jpg')
ON CONFLICT (uuid) DO NOTHING ;

