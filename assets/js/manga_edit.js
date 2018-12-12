window.onload = function() {
	
	var inputValor = document.querySelector("#valor");

	console.log(inputValor);

	var smallPriceText = document.querySelector("#price");

	console.log(smallPriceText)

	smallPriceText.innerHTML = "Preço: " + inputValor.value + ",00 R$";

	$("#valor").keyup(function() {

		var valor = $("#valor").val()

		if (valor != "") {
			$("#price").html( "Preço: " + $("#valor").val() + ",00 R$" );
		} else {
			$("#price").html( "Informe um valor R$: " );
		}
		
	});


	/* AJAX POST VALIDATION */

	var btn_submit = document.querySelector("#btn-submit");

	btn_submit.onclick = function() {

		console.log("formulário submetido");

		const inputHiddenId  = document.querySelector("#id");
		const inputDescricao = document.querySelector("#descricao");

		console.log([ inputHiddenId.value, inputDescricao.value ]);

		ajaxValidation( inputHiddenId.value, inputDescricao.value );

	}

	function ajaxValidation(id, descricao) {

		$.ajax({
			type: 'POST',
			dataType: 'html',
			url: "/manga/edit-validate",
			beforeSend: function() {
				console.log("carregando...");
			},
			data: { id: id,  descricao: descricao },
			success: function(response) {

				console.log(response);

				// Valores esperados: "IsValid" => boolean, "Message" => string
				var json = JSON.parse(response);

				console.log( json.IsValid );

				if ( !json.IsValid ) {

					var htmlMessage = "<p class='alert alert-warning'> Descrição ja existe! </p>";

					$("#validate-message").html( htmlMessage  );

				} else {

					$("#form-manga-edit").submit();

				}
				
			}
		});

	}

} // end onload function
