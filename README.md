## First time service setup:
1. create venv with: 
```
    python -m venv python-env
```
2. activate venv:
```
    wind: python-env\Scripts\activate.bat
    unix: source python-env/Scripts/activate
```
3. install dependicies:
```
    python -m pip install -r requirements.txt
```

## You can run next commands with build.bat:
- dev
```
    Makes docker image of the scraping service sutiable for the dev -
    it has fast build.
```

- go
```
    Create executable of the main go file
    (/cmd/scraing_service/main.go).
```
