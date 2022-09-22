
/******************************************************************************

                      Task 1:  Write a Query
1.Change the website field, so it only contains the domain.
	‣ Example: https://domain.com/index.php → domain.com
2.Count how many spots contain the same domain
3.Return spots which have a domain with a count greater than 1
4.Make a PL/SQL function for point 1 above. [Bonus]

*******************************************************************************/


CREATE OR REPLACE FUNCTION get_rep_domains () RETURNS TABLE (
        domain TEXT
) AS $$
BEGIN
    RETURN QUERY SELECT
        website_domain
    FROM
        (SELECT COUNT(website_domain) AS domain_count, website_domain FROM (SELECT website,
  SUBSTRING("MY_TABLE".website FROM '(?:.*://)?(?:www\.)?([^/?]*)') AS website_domain
FROM "MY_TABLE") AS new_table GROUP BY new_table.website_domain) AS new_table2
    WHERE
        domain_count > 1 ;
END ; $$ LANGUAGE 'plpgsql';

/*
	Call to the function
*/
SELECT * FROM get_rep_domains();
