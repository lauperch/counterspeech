from flask import Flask, render_template, request

app = Flask(__name__)


@app.route('/')
def root():
    return render_template('index.html')

@app.route('/search')
def search():
    searchTerm = request.args.get('q')
    return render_template('index.html', searchTerm=searchTerm)

if __name__ == "__main__":
    app.run(debug=True, host='0.0.0.0')
