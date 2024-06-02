function removeClassesByPrefix(el, prefix) {
  for (var i = el.classList.length - 1; i >= 0; i--) {
    if (el.classList[i].startsWith(prefix)) {
      el.classList.remove(el.classList[i]);
    }
  }
}

document.addEventListener("alpine:init", () => {
  Alpine.data("formValidator", () => ({
    isValid: false,
    inputElements: [],
    init() {
      Iodine.rule(
        "matchingPassword",
        (value) => value === document.getElementById("password").value
      );
      Iodine.setErrorMessage(
        "matchingPassword",
        "Password confirmation needs to match password"
      );
      this.inputElements = [...this.$el.querySelectorAll("input[data-rules]")];
      this.initDomData();
      this.updateErrorMessages();
    },
    initDomData: function () {
      this.inputElements.map((ele) => {
        this[ele.name] = {
          serverErrors: JSON.parse(ele.dataset.serverErrors),
          blurred: false,
        };
      });
    },
    updateErrorMessages: function () {
      this.inputElements.map((ele) => {
        this[ele.name].errorMessage = this.getErrorMessage(ele);
      });
    },
    getErrorMessage: function (ele) {
      this.validate();
      if (this[ele.name].serverErrors.length > 0) {
        return input.serverErrors[0];
      }
      const error = Iodine.assert(ele.value, JSON.parse(ele.dataset.rules));
      if (!error.valid && this[ele.name].blurred) {
        return error.error;
      }
      return "";
    },
    validate: function () {
      const valid = this.inputElements.filter((el) => {
        return Iodine.assert(el.value, JSON.parse(el.dataset.rules)).valid;
      }).length;
      this.isValid = valid === this.inputElements.length;
    },
    change: function (event) {
      if (!this[event.target.name]) {
        return false;
      }
      if (event.type === "input") {
        const formResponse = document.getElementById("form-response");
        removeClassesByPrefix(formResponse, "p-");
        formResponse.innerHTML = "";
        this[event.target.name].serverErrors = [];
      }
      if (event.type === "focusout") {
        this[event.target.name].blurred = true;
      }
      this.updateErrorMessages();
    },
  }));
});
