CREATE TABLE entries (
  id INT AUTO_INCREMENT PRIMARY KEY,
  date DATE,
  country CHAR(2),
  state CHAR(2),
  infected INT,
  hospitalized INT,
  critical INT,
  dead INT,
  recovered INT,
  active INT
)