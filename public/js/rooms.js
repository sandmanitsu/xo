function createRoom() {
    $('.createBtn').on('click', function(evt) {
        rooms = $('.room')

        let id = []
        rooms.each(function(index, item) {
            id.push($(item).attr('id'))
        })

        roomid = Math.max.apply(null, id) + 1;
        if (roomid == -Infinity) {
            roomid = 1
        }
        console.log(roomid);

        $.ajax({
            method: 'post',
            url: '/ws/createRoom',
            contentType: 'application/json; charset=utf-8',
            dataType: 'json',
            data: JSON.stringify({
                'id': String(roomid),
                'name':'room' + roomid
            }),
            success: function(data) {
                $('.rooms').append(htmlRoom(data.id, data.name))
            }
        });
    })
}

function getRooms() {
    $.ajax({
        method: 'get',
        url: '/ws/getRooms',
        success: function(data) {
            data.forEach(function(item) {
                console.log(item);
                $('.rooms').append(htmlRoom(item.id, item.name, item.clientsAmount))
            });
        }
    });
}

// return room html
function htmlRoom(id, name, clientsAmount) {
    clientId = $('h2').attr('clientid')
    username = $('h2').attr('username')

    htmlButton = '<form action="/login/'+ id + '/' + clientId + '/' + username + '" method="get"><button class="room_btn" roomId="' + id + '">Join!</button></form>'
    if (clientsAmount >= 2) {
        htmlButton = '<button class="room_btn">Room is full!</button>'
    }

    return '<div class="room" id="'+ id + '"><div class="room_name">' + name + '</div>' + htmlButton + '</div>'
}

function updateRoomList() {
    $('.updateBtn').on('click', function(evt) {
        $('.room').remove()
        getRooms();
    })
}

document.addEventListener('DOMContentLoaded', function() {
    getRooms()
    createRoom()
    updateRoomList()
}, false);