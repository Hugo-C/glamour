package gold

import (
	"io"

	"github.com/muesli/reflow"
	bf "gopkg.in/russross/blackfriday.v2"
)

type ParagraphElement struct {
}

func (e *ParagraphElement) Render(w io.Writer, node *bf.Node, tr *TermRenderer) error {
	pre := "\n"
	if node.Prev == nil || (node.Parent != nil && node.Parent.Type == bf.Item) {
		pre = ""
	}

	rules := tr.style[Paragraph]
	if rules == nil {
		_, _ = w.Write([]byte(pre))
	} else {
		renderText(w, rules, pre)

		if rules.Prefix != "" {
			renderText(w, rules, rules.Prefix)
		}
	}
	return nil
}

func (e *ParagraphElement) Finish(w io.Writer, node *bf.Node, tr *TermRenderer) error {
	var indent uint
	suffix := ""
	rules := tr.style[Paragraph]
	if rules != nil {
		indent = rules.Indent
		suffix = rules.Suffix
	}
	if suffix != "" {
		renderText(tr.paragraph, rules, suffix)
	}

	iw := &IndentWriter{
		Indent:  indent,
		Forward: w,
	}

	_, err := iw.Write(reflow.ReflowBytes(tr.paragraph.Bytes(), tr.WordWrap-int(indent)))
	tr.paragraph.Reset()
	return err
}