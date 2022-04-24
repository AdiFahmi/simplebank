BEGIN;

INSERT INTO transfers ( from_account_id, to_account_id, ammount )
VALUES
	(
		1,
		2,
	10 
	);
	
INSERT INTO entries (
  account_id,
  ammount
) VALUES (
  1, -10
);

INSERT INTO entries (
  account_id,
  ammount
) VALUES (
  2, 10
);

UPDATE accounts
SET balance = balance - 10
WHERE id = 1;

UPDATE accounts
SET balance = balance + 10
WHERE id = 2;

COMMIT;
