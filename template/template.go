package template

const (

    // compiled using build_template.py
    IndexHtml = `<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>In-Memory File Store</title>
</head>

<body>
    <div>
        <h2>In-Memory File Store</h2>
        <div>
            <h3>Upload File:</h3>
            <input type="file" id="uploadFile">
        </div>
        <button type="button" id="uploadBtn" onclick="upload()">Upload</button>
    </div>

    <div>
        <ul id="ulist">
        </ul>
    </div>

    <script>
        function upload() {
            const uploadBtn = document.getElementById("uploadBtn");
            const uploadFileInput = document.getElementById("uploadFile");
            if (uploadFileInput.files.length === 0) {
                window.alert("Please select a file to upload");
                return;
            }

            const f = uploadFileInput.files[0];
            uploadFileInput.disabled = true;
            uploadBtn.disabled = true;

            fetch('/file/' + f.name, {
                method: "POST",
                body: f,
            })
                .finally(() => {
                    getList();
                    uploadFileInput.value = [];
                    uploadFileInput.disabled = false;
                    uploadBtn.disabled = false;
                })
                .catch((error) => {
                    console.error("Error:", error);
                    window.alert("Failed to upload file");
                });
        }

        async function getList() {
            const ulist = document.getElementById("ulist");
            ulist.innerHTML = "";
            fetch("/file/list", {
                method: "GET",
            })
                .then((response) => response.json())
                .then((result) => {
                    console.log(result);
                    if (result.error) {
                        window.alert(result.msg);
                        return;
                    }

                    for (let p of result.data) {
                        console.log(p)
                        let li = document.createElement("li");
                        let innerLink = document.createElement("a");
                        innerLink.href = 'file/' + p;
                        innerLink.textContent = p;
                        li.appendChild(innerLink);
                        li.setAttribute("target", "_blank");
                        li.style.wordBreak = "break-all";
                        ulist.appendChild(li);
                    }
                })
                .catch((error) => {
                    console.error("Error:", error);
                    window.alert("Failed to fetch file list");
                });
        }

        function getListTask() {
            getList();
            setTimeout(() => {
                getListTask()
            }, 10000);
        }

        getListTask();
    </script>
</body>

</html>`

)

