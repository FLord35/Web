var post = {
    "Title": null,
    "Subtitle": null,
    "AuthorName": null,
    "AuthorPhoto": null,
    "AuthorPhotoBase64": null,
    "Date": null,
    "BigImage": null,
    "BigImageBase64": null,
    "SmallImage": null,
    "SmallImageBase64": null,
    "Content": null
}

var title;
var subtitle;
var authorName;
var date;
var content;

var APhotoB64;
var BImageB64;
var SImageB64;

var FirstAnimationStop=true;

function Publish()
{
    let error = ValidateForm()

    if (error === 0)
    {
        console.log("HEYHEY")
        Preview();
    
        post.Title = title.value;
        post.Subtitle = subtitle.value;
        post.AuthorName = authorName.value;
        post.Date = date.value;
        post.Content = content.value;
        post.AuthorPhoto = document.getElementById("input-author-name").files[0].name;
        post.BigImage = document.getElementById("input-big-image").files[0].name;
        post.SmallImage = document.getElementById("input-small-image").files[0].name;
        post.AuthorPhotoBase64 = APhotoB64;
        post.BigImageBase64 = BImageB64;
        post.SmallImageBase64 = SImageB64;
        
        let XHR = new XMLHttpRequest();
        XHR.open("POST", "/api/post");
        console.log(JSON.stringify(post));
        XHR.send(JSON.stringify(post));
    
        document.getElementById("message").classList.remove('messages__error-message')
        document.getElementById("message").classList.add('messages__success-message')
        document.getElementsByClassName("message__text").item(0).textContent = "Publish Successful!"
        document.getElementsByClassName("message__icon").item(0).src = "../static/img/check-circle.svg"
    }
    else
    {
        document.getElementById("message").classList.remove('messages__success-message')
        document.getElementById("message").classList.add('messages__error-message')
        document.getElementsByClassName("message__text").item(0).textContent = "Whoops! Some fields need your attention :o"
        document.getElementsByClassName("message__icon").item(0).src = "../static/img/alert-circle.svg"
    }

    document.getElementById("message").classList.remove('message-disappear')
    document.getElementById("message").classList.add('message-appear')
    setTimeout(ShowMessage, 100);
    if (FirstAnimationStop)
        FirstAnimationStop = false;
}

function ValidateForm()
{
    let error = 0;
    let formChecker = document.getElementsByClassName("_input")

    for (let index = 0; index < formChecker.length; index++) {
        input = formChecker[index];
        if (input.value === "")
            error++;
    };
    return error;
}

function ShowMessage()
{
    document.getElementById("message").classList.add('message-shown')
    document.getElementById("message").classList.remove('message-hidden')
}

function HideMessage()
{
    document.getElementById("message").classList.add('message-hidden')
    document.getElementById("message").classList.remove('message-shown')
}

function Preview()
{
    title = document.getElementById("titleInput");
    subtitle = document.getElementById("subtitleInput");
    authorName = document.getElementById("authorNameInput");
    date = document.getElementById("dateInput");
    content = document.getElementById("contentInput");

    if (title.value === "")
    {
        document.getElementsByClassName("big-preview__title").item(0).textContent = "New Post";
        document.getElementsByClassName("small-preview__title").item(0).textContent = "New Post";
    }
    else
    {
        document.getElementsByClassName("big-preview__title").item(0).textContent = title.value;
        document.getElementsByClassName("small-preview__title").item(0).textContent = title.value;
    }

    if (subtitle.value === "")
    {
        document.getElementsByClassName("big-preview__subtitle").item(0).textContent = "Please, enter any description";
        document.getElementsByClassName("small-preview__subtitle").item(0).textContent = "Please, enter any description";
    }
    else
    {
        document.getElementsByClassName("big-preview__subtitle").item(0).textContent = subtitle.value;
        document.getElementsByClassName("small-preview__subtitle").item(0).textContent = subtitle.value;
    }

    if (authorName.value === "")
        document.getElementsByClassName("card-info__author-name").item(0).textContent = "Enter author name";
    else
        document.getElementsByClassName("card-info__author-name").item(0).textContent = authorName.value;

    if (date.value === "")
        document.getElementsByClassName("card-info__date").item(0).textContent = "yyyy-mm-dd";    
    else
        document.getElementsByClassName("card-info__date").item(0).textContent = date.value;

    if(!FirstAnimationStop)
        document.getElementById("message").classList.add('message-disappear');
        
    document.getElementById("message").classList.remove('message-appear')

    setTimeout(HideMessage, 100);
}

//Наводим красоту для входящих фотографий

function PreviewAuthorPhoto()
{
    if(document.getElementById("input-author-name").files.length == 0) //Если нет загруженного изображения
    {
        document.getElementsByClassName("photo-label__text").item(0).textContent = "Upload";
        document.getElementsByClassName("photo-label__icon").item(0).style.width = "0px";
        document.getElementsByClassName("author-photo-frame__remove").item(0).style.visibility = "hidden";
        document.getElementsByClassName("author-photo-frame__preview").item(0).style.display = "none";
        document.getElementsByClassName("photo-label__photo").item(0).style.display = "flex";
        document.getElementsByClassName("card-info__author-photo").item(0).src= "../static/img/author_photo_placeholder.svg";
    }
    else
    {
        document.getElementsByClassName("photo-label__text").item(0).textContent = "Upload New";
        document.getElementsByClassName("photo-label__icon").item(0).style.width = "34px";
        document.getElementsByClassName("author-photo-frame__remove").item(0).style.visibility = "visible";
        document.getElementsByClassName("author-photo-frame__preview").item(0).style.display = "flex";
        document.getElementsByClassName("photo-label__photo").item(0).style.display = "none";

        const preview = document.querySelector("div.author-photo__author-photo-frame img"); //author-photo-frame__preview
        const previewRight = document.querySelector("div.small-preview__card-info img"); //card-info__author-photo
        const file = document.querySelector("div.author-photo__author-photo-frame input[type=file]").files[0];
        const reader = new FileReader();

        reader.addEventListener(
            "load",
            () => {
              preview.src = reader.result;
              previewRight.src = reader.result;
              APhotoB64 = reader.result;
            },
            false,
          );

        if (file) {
            reader.readAsDataURL(file);
        }
    }
}

function RemoveAuthorPhoto()
{
    document.getElementById("input-author-name").value = "";
    document.querySelector("div.author-photo__author-photo-frame img").src = "";
    PreviewAuthorPhoto();
}

function PreviewBigImage() 
{
    if(document.getElementById("input-big-image").files.length == 0) //Если нет загруженного изображения
    {
        document.getElementsByClassName("big-image-content__big-bonus-menu").item(0).style.display = "none";
        document.getElementsByClassName("big-image-content__text").item(0).style.display = "inherit";
        document.getElementsByClassName("big-image__preview").item(0).style.display = "none";
        document.getElementsByClassName("big-image__big-image-label").item(0).style.display = "flex";
        document.getElementsByClassName("big-preview__subimage").item(0).style.display = "none";
    }
    else  
    {
        document.getElementsByClassName("big-image-content__big-bonus-menu").item(0).style.display = "flex";
        document.getElementsByClassName("big-image-content__text").item(0).style.display = "none";
        document.getElementsByClassName("big-image__preview").item(0).style.display = "flex";
        document.getElementsByClassName("big-image__big-image-label").item(0).style.display = "none";
        document.getElementsByClassName("big-preview__subimage").item(0).style.display = "flex";

        const preview = document.querySelector("div.big-image-content__big-image img"); //big-image__preview
        const previewRight = document.querySelector("div.big-preview__image img"); //big-preview__subimage
        const file = document.querySelector("div.big-image-content__big-image input[type=file]").files[0];
        const reader = new FileReader();

        reader.addEventListener(
            "load",
            () => {
              preview.src = reader.result;
              previewRight.src = reader.result;
              BImageB64 = reader.result;
            },
            false,
          );

        if (file) {
            reader.readAsDataURL(file);
        }
    }
}

function RemoveBigPhoto()
{
    document.getElementById("input-big-image").value = "";
    document.querySelector("div.big-image-content__big-image img").src = "";
    PreviewBigImage();
}

function PreviewSmallImage()
{
    if(document.getElementById("input-small-image").files.length == 0) //Если нет загруженного изображения
    {
        document.getElementsByClassName("small-image-content__small-bonus-menu").item(0).style.display = "none";
        document.getElementsByClassName("small-image-content__text").item(0).style.display = "inherit";
        document.getElementsByClassName("small-image__preview").item(0).style.display = "none";
        document.getElementsByClassName("small-image__small-image-label").item(0).style.display = "flex";
        document.getElementsByClassName("small-preview__subimage").item(0).style.display = "none";
    }
    else
    {
        document.getElementsByClassName("small-image-content__small-bonus-menu").item(0).style.display = "flex";
        document.getElementsByClassName("small-image-content__text").item(0).style.display = "none";
        document.getElementsByClassName("small-image__preview").item(0).style.display = "flex";
        document.getElementsByClassName("small-image__small-image-label").item(0).style.display = "none";
        document.getElementsByClassName("small-preview__subimage").item(0).style.display = "flex";

        const preview = document.querySelector("div.small-image-content__small-image img"); //small-image__preview
        const previewRight = document.querySelector("div.small-preview__image img"); //big-preview__subimage
        const file = document.querySelector("div.small-image-content__small-image input[type=file]").files[0];
        const reader = new FileReader();

        reader.addEventListener(
            "load",
            () => {
              preview.src = reader.result;
              previewRight.src = reader.result;
              SImageB64 = reader.result;
            },
            false,
          );

        if (file) {
            reader.readAsDataURL(file);
        }
    }
}

function RemoveSmallPhoto()
{
    document.getElementById("input-small-image").value = "";
    document.querySelector("div.small-image-content__small-image img").src = "";
    PreviewSmallImage();
}