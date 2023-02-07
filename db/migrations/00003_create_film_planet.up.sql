BEGIN;

CREATE TABLE film_planet (
  planet_id INT NOT NULL,
  film_id INT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMP NULL DEFAULT NULL,
  CONSTRAINT planet_fk
    FOREIGN KEY (planet_id)
    REFERENCES planet (id)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT film_fk
    FOREIGN KEY (film_id)
    REFERENCES film (id)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION);  

COMMIT;