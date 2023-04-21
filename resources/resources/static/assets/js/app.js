(function ($) {
  $(function () {
    /**
     * タスク登録／更新ダイアログを表示する
     */
    var dialog = $('#task-form').dialog({
      autoOpen: false,
      width: 500,
      height: 550,
      modal: true,
      buttons: {
        Ok: function () {
          updateDetail(
            $('#taskid').val(),
            $('#tasktype').val(),
            $('#taskname').val(),
            $('#plan').val(),
            $('#result').val(),
            $('#unit').val(),
            $('#due').val(),
            $('#person').val()
          )
            .done(function (data) {
              if (data.result) {
                $.growl.notice({title: '', message: data.message});
                dialog.dialog('close');
                getTask().
                  done(function (data) {
                    if (data.result) {
                      showTask(data);
                    }
                  });
              }
              else {
                $.growl.warning({title: '', message: data.message});
              }
            })
            .fail(function (data) {
              $.growl.error({message: 'failed.'});
            });
        },
        Cancel: function () {
          dialog.dialog('close');
        }
      },
      close: function () {
        form[0].reset();
      }
    }),
        form = dialog.find('form').on('submit', function (event) {
          event.preventDefault();
          updateDetail(
            $('#taskid').val(),
            $('#tasktype').val(),
            $('#taskname').val(),
            $('#plan').val(),
            $('#result').val(),
            $('#unit').val(),
            $('#due').val(),
            $('#person').val()
          )
            .done(function (data) {
              if (data.result) {
                $.growl.notice({title: '', message: data.message});
                dialog.dialog('close');
                getTask().
                  done(function (data) {
                    if (data.result) {
                      showTask(data);
                    }
                  });
              }
              else {
                $.growl.warning({title: '', message: data.message});
              }
            })
            .fail(function (data) {
              $.growl.error({title: '', message: 'failed.'});
            });
        });

    $.datepicker.setDefaults($.datepicker.regional['ja']);
    $('#datepicker').datepicker();
    
    $('#create-task').button().on('click', function () {
      $('#taskid').val("");
      dialog.dialog('open');
    });

    $('#refresh-task').button().on('click', function () {
      getTask().
        done(function (data) {
          if (data.result) {
            showTask(data);
            $.growl.notice({title: '', message: 'updated'});
          }
          else {
            $.growl.warning({title: '', message: 'reload error'});
          }
        });
    });
    
    /*-
     * タスクを傾ける
     */
    function tilt_direction (item) {
      var left_pos = item.position().left,
          move_handler = function (e) {
            if (e.pageX >= left_pos) {
              item.addClass('right');
              item.removeClass('left');
            }
            else {
              item.addClass('left');
              item.removeClass('right');
            }
            left_pos = e.pageX;
          };
      $('html').bind('mousemove', move_handler);
      item.data('move_handler', move_handler);
    }
    
    /*-
     * 各レーンを移動できるようにする
     */
    $('.list-cards').sortable({
      connectWith: '.list-cards',
      start: function (event, ui) {
        ui.item.addClass('tilt');
        tilt_direction(ui.item);
      },
      stop: function (event, ui) {
        ui.item.removeClass('tilt');
        $('html').unbind('mousemove', ui.item.data('move_handler'));
        ui.item.removeData('move_handler');
      },
      receive: function (event, ui) {
        var item = ui.item.context.id,
            from = ui.sender.context.id,
            to = $(this).context.id;

        updateStatus(item, from, to)
          .done(function (data) {
          })
          .fail(function () {
            $.growl.error({message: 'server error'});
          });
      }
    });

    /**
     * アプリ名称取得
     *
     */
    function getApp() {
      var res = $.Deferred();
      $.ajax({
        method: 'GET',
        url: '/api/app/name',
        success: res.resolve,
        error: res.reject
      });

      return res.promise();
    }

    /**
     * ユーザ一覧を取得する
     *
     */
    function getUser() {
      var res = $.Deferred();
      $.ajax({
        method: 'GET',
        url: '/api/user',
        success: res.resolve,
        error: res.reject
      });

      return res.promise();
    }
    
    /**
     * タスク詳細の追加、更新を行う
     *
     */
    function updateDetail (id, type, name, plan, result, unit, due, person) {
      var res = $.Deferred(),
          method,
          data,
          url;
      
      if (!id) {
        method = 'POST';
        data = {
          type: type,
          name: name,
          plan: plan,
          result: result,
          unit: unit,
          due: due,
          person: person
        };
        url = '/api/task';
        
      }
      else {
        method = 'PUT';
        data = {
          id: id,
          type: type,
          name: name,
          plan: plan,
          result: result,
          unit: unit,
          due: due,
          person: person
        };
        url = '/api/task/' + id;
      }
      
      $.ajax({
        method: method,
        contentType: 'application/json;charset=UTF-8',
        dataType: 'json',
        url: url,
        data: JSON.stringify(data),
        success: res.resolve,
        error: res.reject
      });
      return res.promise();
    }
    
    /**
     * タスクのステータスを更新する
     *
     */
    function updateStatus (id, from, to) {
      var res = $.Deferred(),
          data = {
            status: to
          };
      $.ajax({
        method: 'PUT',
        url: '/api/task/' + id,
        data: JSON.stringify(data),
        success: res.resolve,
        error: res.reject
      });
      return res.promise();
    }
    
    /**
     * タスク一覧を取得する
     *
     */
    function getTask () {
      var res = $.Deferred(),
          user = $('#search-user').val(),
          update = $('#datepicker').val();

      $.ajax({
        method: 'GET',
        url: '/api/search?user=' + user + '&update=' + update,
        success: res.resolve,
        error: res.reject
      });
      return res.promise();
    }

    /**
     * タスクのDOMを構築する
     *
     */
    function showTask (data) {
      $('#task-todo').empty();
      $('#task-inprogress').empty();
      $('#task-done').empty();
      $('#task-archive').empty();

      var func = function (i, val) {
        var task = $('<div></div>',
                     {
                       'class': 'list-card js-card',
                       'id': val.id,
                       on: {
                         /**
                          * タスクをクリックした際に、編集用ダイアログを表示する。
                          */
                         click: function (event) {
                           $('#taskid').val($('.column-taskid', this).text());
                           $('#tasktype').val($('.column-tasktype', this).text());
                           $('#taskname').val($('.column-taskname', this).text());
                           $('#plan').val($('.column-plan', this).text());
                           $('#result').val($('.column-result', this).text());
                           $('#unit').val($('.column-unit', this).text());
                           $('#due').val($('.column-due', this).text());
                           $('#person').val($('.column-person', this).text());
    
                           dialog.dialog('open');
                         }
                       }
                     }),
            desc = $('<div></div>',
                     {
                       'style': 'display: none;'
                     });


        desc.append($('<div></div>', {'class': 'column-taskid'}).append(val.id))
          .append($('<div></div>', {'class': 'column-tasktype'}).append(val.type))
          .append($('<div></div>', {'class': 'column-taskname'}).append(val.name))
          .append($('<div></div>', {'class': 'column-plan'}).append(val.plan))
          .append($('<div></div>', {'class': 'column-result'}).append(val.result))
          .append($('<div></div>', {'class': 'column-unit'}).append(val.unit))
          .append($('<div></div>', {'class': 'column-due'}).append(val.due))
          .append($('<div></div>', {'class': 'column-person'}).append(val.person))
        ;

        task.append(desc)
          .append('<div class="title-header">' + val.type + ' : ' + val.name + '</div>' +
                  '<div align="right">' + val.result + ' / ' + val.plan + ' ' + val.unit + '</div>' +
                  '<div align="right">Due:' + val.due + '</div>' +
                  '<div align="right">by ' + val.person + '</div>');

        return task;
      };

      $.each(data.tasktodo, function (i, val) {
        $('#task-todo').append(func(i, val));
      });

      $.each(data.taskinprogress, function (i, val) {
        $('#task-inprogress').append(func(i, val));
      });

      $.each(data.taskdone, function (i, val) {
        $('#task-done').append(func(i, val));
      });

      $.each(data.taskarchive, function (i, val) {
        $('#task-archive').append(func(i, val));
      });
    }

    getApp().
      done(function (data) {
        $('#appname').text(data.appname);
      });

    getUser().
      done(function (data) {
        if (data.result) {
          $('#search-user').empty();
          $('#search-user').append('<option value="" selected></option>');
          $.each(data.users, function (i, val) {
            $('#search-user').append('<option value="' + val.username + '">' + val.username + '</option>');
          });
        }
      });
    
    getTask().
      done(function (data) {
        if (data.result) {
          showTask(data);
        }
      });
  });
  
})(window.$)
;
  
  
