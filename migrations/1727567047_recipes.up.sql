
CREATE TABLE IF NOT EXISTS recipes (
  id UUID PRIMARY KEY,
  name VARCHAR (100) NOT NULL,
  description TEXT
);

CREATE TABLE IF NOT EXISTS recipe_ingredients (
  recipe_id UUID REFERENCES recipes(id),
  ingredient_id UUID REFERENCES ingredients(id),
  PRIMARY KEY (recipe_id, ingredient_id)
);
