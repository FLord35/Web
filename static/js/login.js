var FirstAnimationStop=true;

function LogIn()
{
    let error = ValidateUser()

    if (error === 0)
    {
        alert("Молодец")
    }
    else
    {
        document.getElementById("message").classList.add('messages__error-message')
        document.getElementsByClassName("message__text").item(0).textContent = "A-ah! Check all fields"
        document.getElementById("message").classList.remove('message-disappear')
        document.getElementById("message").classList.add('message-appear')
        setTimeout(LogInShowMessage, 100);
        if (FirstAnimationStop)
        FirstAnimationStop = false;
    }

}

function ValidateUser()
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

function LogInShowMessage()
{
    document.getElementById("message").classList.add('message-shown')
    document.getElementById("message").classList.remove('message-hidden')
}

function LogInHideMessage()
{
    document.getElementById("message").classList.add('message-hidden')
    document.getElementById("message").classList.remove('message-shown')
}

function ChangeFieldColor()
{
    fieldEmail = document.getElementById("input-email");
    fieldPassword = document.getElementById("input-password");

    if (fieldEmail.value === "")
    {
        fieldEmail.style.background = "#FFFFFF";
        fieldEmail.style.borderbottom = "1px solid #EAEAEA";
    }
    else
    {
        fieldEmail.style.background = "#F7F7F7";
        fieldEmail.style.borderBottom = "1px solid #2E2E2E";
    }

    if (fieldPassword.value === "")
    {
        fieldPassword.style.background = "#FFFFFF";
        fieldPassword.style.borderbottom = "1px solid #EAEAEA";
    }
    else
    {
        fieldPassword.style.background = "#F7F7F7";
        fieldPassword.style.borderBottom = "1px solid #2E2E2E";
    }

    if(!FirstAnimationStop)
        document.getElementById("message").classList.add('message-disappear')
    document.getElementById("message").classList.remove('message-appear')
    setTimeout(LogInHideMessage, 100);
}

function ChangePasswordVisibility()
{
    passwordField = document.getElementById("input-password");
    if (passwordField.type === "password")
    {
        document.getElementById("input-password").type = "text";
        document.getElementsByClassName("email-password__image").item(0).src = "../static/img/eye.svg"
    }
    else
    {
        document.getElementById("input-password").type = "password";
        document.getElementsByClassName("email-password__image").item(0).src = "../static/img/eye-off.svg"
    }
}