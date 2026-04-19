import pandas as pd
from sqlalchemy import create_engine, text

# 1. Configurações de Conexão (Ajuste com seus dados)
DB_USER = 'postgres'
DB_PASSWORD = 'postgres'
DB_HOST = 'postgres'
DB_PORT = '5432'
DB_NAME = 'postgres'

engine = create_engine(f'postgresql://{DB_USER}:{DB_PASSWORD}@{DB_HOST}:{DB_PORT}/{DB_NAME}')

def create_tables():
    """Cria as tabelas no banco de dados caso não existam"""
    sql_statements = [
        """
        CREATE TABLE IF NOT EXISTS movies (
            movie_id VARCHAR(20) PRIMARY KEY,
            title TEXT,
            year INTEGER,
            genre VARCHAR(100),
            director VARCHAR(255)
        );
        """,
        """
        CREATE TABLE IF NOT EXISTS users (
            user_id INTEGER PRIMARY KEY,
            name VARCHAR(255),
            email VARCHAR(255),
            country VARCHAR(100)
        );
        """,
        """
        CREATE TABLE IF NOT EXISTS ratings (
            user_id INTEGER REFERENCES users(user_id),
            movie_id VARCHAR(20) REFERENCES movies(movie_id),
            rating NUMERIC(3,1),
            timestamp DATE
        );
        """
    ]
    
    with engine.connect() as conn:
        print("Creating tables")
        for statement in sql_statements:
            conn.execute(text(statement))
        conn.commit()
        print("Tables created")

def upload_csv_to_postgres(file_path, table_name):
    try:
        print(f"reading {file_path}...")
        df = pd.read_csv(file_path)
        
        if 'timestamp' in df.columns:
            df['timestamp'] = pd.to_datetime(df['timestamp'])

        df.to_sql(table_name, engine, if_exists='append', index=False)
        
        print(f"Success: {len(df)} entries imported to '{table_name}'.")
    except Exception as e:
        print(f"Error at {table_name}: {e}")

if __name__ == "__main__":
    create_tables()
    files_to_upload = {
        'movies.csv': 'movies',
        'users.csv': 'users',
        'ratings.csv': 'ratings'
    }

    for file, table in files_to_upload.items():
        upload_csv_to_postgres(file, table)