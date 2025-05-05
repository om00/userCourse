CREATE TABLE courses (
    id INT AUTO_INCREMENT PRIMARY KEY,  
    name VARCHAR(100) NOT NULL,        
    description TEXT,                  
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP 
);