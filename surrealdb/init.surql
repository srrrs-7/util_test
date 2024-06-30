-- Create a schemafull table
DEFINE TABLE user SCHEMAFULL;

-- Specify fields on the user table
DEFINE FIELD name ON TABLE user TYPE object;
DEFINE FIELD name.first ON TABLE user TYPE string;
DEFINE FIELD name.last ON TABLE user TYPE string;
DEFINE FIELD email ON TABLE user TYPE string ASSERT string::is::email($value);

-- Add a unique index on the email field preventing duplicate values
DEFINE INDEX email ON TABLE user COLUMNS email UNIQUE;

-- Create a new event whenever a user changes their email address
DEFINE EVENT email ON TABLE user WHEN $before.email != $after.email THEN (
    CREATE event SET user = $value, time = time::now(), value = $after.email, action = 'email_changed'
);
