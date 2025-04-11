DROP TABLE IF EXISTS log;
CREATE TABLE log (
    id INT AUTO_INCREMENT NOT NULL PRIMARY KEY,
    workout_name VARCHAR(128) NOT NULL,
    comment VARCHAR(255),
    date DATE NOT NULL,
    number_of_sets INT NOT NULL,
    number_of_reps INT NOT NULL,
    weight INT NOT NULL,
    effort INT NOT NULL
);

INSERT INTO log (workout_name, comment, date, number_of_sets, number_of_reps, weight, effort) 
VALUES 
    ('Back & biceps', 'felts strong', '2025-04-10', 3, 12, 100, 5),
    ('Chest & triceps', 'felts strong', '2025-04-10', 3, 12, 100, 5),
    ('Leg day', 'felts strong', '2025-04-10', 3, 12, 160, 5),
    ('Shoulder work', 'felts weak', '2025-04-10', 3, 12, 60, 5);
