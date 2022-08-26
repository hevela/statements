### file names
Place here the files to parse, file names must be a valid email with an `.csv` extension

*example*
`vell.once@gmail.com.csv`

### CSV structure
the CSV file contents must follow the following structure:

| ID  | DATE       | Transaction  |
|-----|------------|--------------|
| int | YYYY-MM-DD | signed float |

*example:*

```
ID,Date,Transaction
0,2021-07-15,+60.5
1,2021-07-28,-10.3
2,2021-08-02,-20.46
3,2021-08-13,+10
```
## Notes
* The first row must be the headers of the CSV
* All amounts must be prepended by a + or - sign

