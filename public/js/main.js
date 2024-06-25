function ActionListener() {
    var data = {
        row1: {
            a: '',
            b: '',
            c: '',
        },
        row2: {
            a: '',
            b: '',
            c: '',
        },
        row3: {
            a: '',
            b: '',
            c: '',
        },
    }

    cells = $('.column')
    var turn = 'X'

    cells.each(function(index) {
        $(this).on('click', function () {
            className = $(this).attr('class')
            row = className.slice(18,22);
            cell = className.slice(23);

            data[row][cell] = turn;

            if ($(this).text() != 'N') {
                return
            }

            $(this).html(turn)
            $(this).css('color','black')

            if (turn == 'X') {
                turn = 'O';
            } else {
                turn = 'X'
            }

            $('.turn').html(turn)
        })
    })
}

document.addEventListener('DOMContentLoaded', function() {
    ActionListener()
}, false);