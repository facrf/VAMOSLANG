package texto_util

import (
	"reflect"
	"testing"
)

func TestTextoUtil(t *testing.T) {
	// Dividir e Juntar
	partes := Dividir("VAMOS,LANG,Go,Compilador", ",")
	esperadoPartes := []string{"VAMOS", "LANG", "Go", "Compilador"}
	if !reflect.DeepEqual(partes, esperadoPartes) {
		t.Errorf("Dividir falhou: obteve %v, esperava %v", partes, esperadoPartes)
	}

	junto := Juntar(partes, " - ")
	if junto != "VAMOS - LANG - Go - Compilador" {
		t.Errorf("Juntar falhou: obteve '%s'", junto)
	}

	// Contem, ComecaCom, TerminaCom
	frase := "VAMOS-LANG é a linguagem do futuro"
	if !Contem(frase, "linguagem") {
		t.Errorf("Contem falhou para 'linguagem'")
	}
	if Contem(frase, "Python") {
		t.Errorf("Contem não deveria encontrar 'Python'")
	}
	if !ComecaCom(frase, "VAMOS") {
		t.Errorf("ComecaCom falhou para prefixo 'VAMOS'")
	}
	if !TerminaCom(frase, "futuro") {
		t.Errorf("TerminaCom falhou para sufixo 'futuro'")
	}

	// Maiúsculas, Minúsculas e Aparar
	textoEspacos := "   Olá Mundo VAMOS   "
	if Aparar(textoEspacos) != "Olá Mundo VAMOS" {
		t.Errorf("Aparar falhou: obteve '%s'", Aparar(textoEspacos))
	}
	if ParaMaiusculas("vamos") != "VAMOS" {
		t.Errorf("ParaMaiusculas falhou")
	}
	if ParaMinusculas("VAMOS") != "vamos" {
		t.Errorf("ParaMinusculas falhou")
	}

	// Substituir
	subst := Substituir("banana maçã banana", "banana", "laranja")
	if subst != "laranja maçã laranja" {
		t.Errorf("Substituir falhou: obteve '%s'", subst)
	}

	// Repetir
	if Repetir("AB", 3) != "ABABAB" {
		t.Errorf("Repetir falhou: obteve '%s'", Repetir("AB", 3))
	}
	if Repetir("AB", 0) != "" {
		t.Errorf("Repetir com zero deveria ser vazio")
	}

	// IndiceDe
	if IndiceDe("programação", "gra") != 3 {
		t.Errorf("IndiceDe falhou: esperado 3, obteve %d", IndiceDe("programação", "gra"))
	}
	if IndiceDe("programação", "xyz") != -1 {
		t.Errorf("IndiceDe inexistente deveria ser -1")
	}

	// Inverter (suporte a UTF-8 / acentos)
	invertido := Inverter("VAMOS")
	if invertido != "SOMAV" {
		t.Errorf("Inverter falhou: obteve '%s'", invertido)
	}
	invertidoAcento := Inverter("café")
	if invertidoAcento != "éfac" {
		t.Errorf("Inverter com acento falhou: obteve '%s'", invertidoAcento)
	}
}
