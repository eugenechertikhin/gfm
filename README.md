# File Manager

Another one try to develop midnight commander clone in golang.

## Command line parameters

Program can start with default parameters or some can be specified in programm arguments

    -a switch to ascii border mode (default is utf8)
    -scheme <scheme> color scheme, one of colour, bw, custom
    -e <filename> star program in editor mode with filename
    -v <filename> star program in viewer mode with filename
    -b <filename> star program in hex editor mode with filename
    -w show license information

## TODO

Many changes and idea still in implementation. Nearest plan below:

- history (commands, mkdir, etc)
- check ESC key, implemnt ESC+Enter, ESC+1, ESC+2, etc
- implement viewer and editor
- implement top menubar
- implement configuration screen
- help window with all shortcuts
- copy, move, create dir 
- menu

## Licensing

Licensed under GNU GENERAL PUBLIC LICENSE version 3

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with this program.  If not, see <https://www.gnu.org/licenses/>.