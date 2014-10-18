(function(w) {
  var Quiz = w.Quiz = (w.Quiz || {});

  var $submit = document.querySelector('.submit-quiz'),
      $questions = document.querySelectorAll('.question'),
      $answers = document.querySelectorAll('.answer'),
      answers = getAnswerMap($answers),
      formData = Quiz.formData = {},
      EMPTY = "No Title Selected",
      state = {
        currentQuestion: null
      };

  Quiz.formData = function() { return formData; };
  Quiz.answers = function() { return answers; };

  function getAnswerMap(elems) {
    var o = {};
    for (var i = 0; i < elems.length; i++) {
      o[elems[i].dataset.id] = elems[i].textContent;
    }
    return o;
  }

  function answerQuestion(questionElem, answerElem) {
    var questionId = questionElem.dataset.id,
        answerId = answerElem.dataset.id;

    formData[questionId] = answerId;
    questionElem.querySelector('figcaption').textContent = answers[questionId];
  }

  function clearAnswer(questionElem) {
    var questionId = questionElem.dataset.id,
        answerId = formData[questionId];

    delete formData[questionId];
    questionElem.querySelector('figcaption').textContent = EMPTY;
  }

  function setupHandlers() {
    $question.addEventListener('click', function(e) {

    });
    $answers.addEventListener('click', function(e) {

    });
  }

})(this);
