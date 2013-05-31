{{template "header"}}

<div class="page">
    <div class="user">
        <div> Hello {{.Name}} </div>
    </div>


    <div>

        <div class="content">
            <div id="tasks"></div>

            <div>
                <form action="/tasks">
                    title
                    <div><input type="text" id="task_title"></div>

                    content
                    <div> <textarea id="task_content"></textarea></div>

                    <div><input type="button" value="add task" id="add" /></div>
                </form>
            </div>
        </div>




    </div>

</div>


<script type="text/javascript">

    var template1 , template2

    var refresh = function(data) {
        var s = template1(data)
        $("#tasks").html(s)
    }

    var renderTask = function(data) {
        return template2(data)
    }

    $(function() {

        template1 = _.template($("#t1").html())
        template2 = _.template($("#t2").html())

        $.getJSON("/tasks", function(data) {
            refresh({tasks: data})

        })

        
        $(".delete").live("click", function(e) {
            var d = $(this)
            var url = d.attr("href")

            $.ajax({
                url: url,
                type: "delete",
                success: function(data) { 
                    d.parent().remove()
                }
            })
            return false
        })

        $("#add").click(function() {
            var title = $("#task_title").val()
            var content = $("#task_content").val()

            $.ajax({
                url: "/tasks",
                type: "post",
                data: { title: title, content: content },
                success: function(data) { 
                    var s = renderTask({task: data})

                    $("#list").append(s)
                }
            })
        })

    });
</script>


<script type="text/template" id="t1">
    
    <ul id="list">
    <% _.each(tasks, function(t){ %>
      
      <%= renderTask({task:t}) %>
   
    <% }); %>
    </ul>
    
</script>

<script type="text/template" id="t2">
      <li>
        <span><%= task.Title %></span>
        <span><%= task.CreatedAt %></span> 
        <a href="/tasks/<%= task.Id %>" class="view">view</a>
        <a href="/tasks/<%= task.Id %>" class="delete">delete</a>
      </li>
    
</script>



{{template "footer"}}
