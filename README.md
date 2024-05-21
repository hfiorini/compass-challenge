# compass-challenge

Instructions:
You have been asked to take a list of contact information and identify which contacts are potentially duplicates.  You need to write a function that will do the following:

1.	Identify which contacts are possible matches. 
2.	A score for each match that represents how accurate you believe the match is.  This scoring is defined by you.
3.	A contact might have multiple or no matches
4.	All processing should be done in working memory (no database).

Example:

Input

Contact ID	First Name	Last Name	Email Addres	Zip Code	Address
1001	C	F	mollis.lectus.pede@outlook.net		449-6990 Tellus. Rd.
1002	C	French	mollis.lectus.pede@outlook.net	39746	449-6990 Tellus. Rd.
1003	Ciara	F	non.lacinia.at@zoho.ca	39746	

Output
ContactID Source	ContactID Match	Accuracy
1001	1002	High
1001	1003	Low

The input file is included in the assessment packet.

As part of your solution you will need to produce:
1.	The code used to identify the potential duplicates
2.	A set of tests that you use to verify the functions
3.	Comments in the code are welcome

