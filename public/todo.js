(function ($) {
    'use strict';
    $(function () {
        var todoListItem = $('.todo-list');
        var todoListInput = $('.todo-list-input');
        $('.todo-list-add-btn').on("click", function (event) {
            event.preventDefault();

            var item = $(this).prevAll('.todo-list-input').val();

            if (item) {
                $.post('/todos', {name: item}, addItem);
                // todoListItem.append("<li><div class='form-check'><label class='form-check-label'><input class='checkbox' type='checkbox' />" + item + "<i class='input-helper'></i></label></div><i class='remove mdi mdi-close-circle-outline'></i></li>");
                todoListInput.val("");
            }

        });

        let addItem = function (item) {
            if (item.completed === true) {
                todoListItem.append("<li class='completed'"+" id='"+item.id+"'><div class='form-check'><label class='form-check-label'><input class='checkbox' checked='checked' type='checkbox' />" + item.name + "<i class='input-helper'></i></label></div><i class='remove mdi mdi-close-circle-outline'></i></li>");
            } else {
                todoListItem.append("<li"+ " id='"+item.id+"'><div class='form-check'><label class='form-check-label'><input class='checkbox' type='checkbox' />" + item.name + "<i class='input-helper'></i></label></div><i class='remove mdi mdi-close-circle-outline'></i></li>");
            }
        };

        $.get('/todos', function (items) {
            items.forEach(e => {
                addItem(e)
            });
        });

        todoListItem.on('change', '.checkbox', function () {
            const id = $(this).closest("li").attr('id');
            const $self = $(this);
            const reqCompleteValue = !$(this).attr('checked');
            $.get("/todoComplete/" + id +"?completed="+reqCompleteValue, function (data) {
                if (reqCompleteValue) {
                    $self.attr('checked', 'checked');
                }else {
                    $self.removeAttr('checked');
                }
                $self.closest("li").toggleClass('completed');
            });
        });

        todoListItem.on('click', '.remove', function () {
            //url = todos/id , method = delete
            const id = $(this).closest("li").attr('id');
            const $self = $(this);
            $.ajax({
                url: 'todos/'+id,
                type: 'DELETE',
                success: function (data) {
                    if (data.success) {
                        $self.parent().remove();
                    } else {
                        alert('실패');
                    }
                },
            })

        });

    });
})(jQuery);