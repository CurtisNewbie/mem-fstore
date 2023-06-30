with open('index.html', 'r') as f:
    index_html = f.read()
    # print(index_html)

    with open('template/template.go', 'w') as fout:
        fout.write(f'''package template

const (

    // compiled using build_template.py
    IndexHtml = `{index_html}`

)

''')

