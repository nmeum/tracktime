# timetrack

Small utility for tracking working hours in a plain text file.

## Usage

Tracking working hours is complicated as one usually has to consider
holidays and vacations for calculating required working hours. For this
reason, this tool takes a different approach. The tool operates under
the assumption that one has to work the same amount of hours each day.
To compensate overtime, the tool additionally tracks a delta between
tracked workdays.

For instance, if one has to work 8 hours each day and worked 9 hours on
Monday only 7 hours would be required on Tuesday:

	$ cat ~/.timetrack
	06.01.2020	0900	1800	Important work stuff
	07.01.2020	1000	1700	More important work stuff
	$ timetrack -h 8 ~/.timetrack
	06.01.2020      9h0m0s  | 1h0m0s
	07.01.2020      7h0m0s  | 0s

The last column of the `timetrack` output represents the delta. As
illustrated, it is zero on the 7th of January even though only 7 hours
of work were done. The overtime done on the 6th of January is used to
compensate for this.

## Input format

The input format uses four tab-separated fields. The fields have the
following meaning:

1. Workday. Multiple entries for the same day are allowed. The
   utilized date format can be customized using the `TIMETRACK_FORMAT`
   environment variable. Refer to the documentation of the
   [go time pkg](https://golang.org/pkg/time/#pkg-constants) for more
   information.
2. Starting time of the described activity without a hours/minutes
   separator and padded with zeros to four digits.
3. End time of the described activity.
4. Description of the activity.

## Installation

No dependencies, simply use `go get`:

	$ go get github.com/nmeum/timetrack

## License

This program is free software: you can redistribute it and/or modify it
under the terms of the GNU General Public License as published by the
Free Software Foundation, either version 3 of the License, or (at your
option) any later version.

This program is distributed in the hope that it will be useful, but
WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General
Public License for more details.

You should have received a copy of the GNU General Public License along
with this program. If not, see <http://www.gnu.org/licenses/>.
