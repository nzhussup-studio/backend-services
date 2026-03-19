TRUNCATE TABLE
  certificates,
  education,
  projects,
  skills,
  work_experience,
  users
RESTART IDENTITY CASCADE;

INSERT INTO users (username, password, role) VALUES
  (
    'admin',
    '$2y$12$uWgXdVGJHlhCbb1jzptSZ.aoICVeMb6Z551HAF8jNBh3zf8Xkyp.C',
    'ROLE_ADMIN'
  ),
  (
    'jane.user',
    '$2y$12$azF2crHGiAjd2CPyfDuuuuF450yqpkRRUsfdqWA7yQCNvsewAzjK.',
    'ROLE_USER'
  );

INSERT INTO work_experience (company, location, start_date, end_date, position, description, display_order, tech_stack) VALUES
  (
    'Freelance / Personal Platform',
    'Remote',
    'Jan 2025',
    NULL,
    'Platform Engineer',
    'Owns the architecture and delivery of a personal cloud platform covering backend services, admin tooling, CI/CD, and Kubernetes-based operations.',
    1,
    'Go,Java,Spring Boot,PostgreSQL,Redis,Docker,Kubernetes,GitHub Actions'
  ),
  (
    'Independent Consulting',
    'Vienna, Austria',
    'May 2024',
    'Dec 2024',
    'Software Engineer',
    'Delivered small-to-medium backend and frontend systems with a focus on maintainable APIs, release automation, and deployment hardening.',
    2,
    'TypeScript,React,Spring Boot,REST,PostgreSQL,Docker'
  ),
  (
    'Media Systems Lab',
    'Munich, Germany',
    'Oct 2022',
    'Apr 2024',
    'Backend Developer',
    'Built internal APIs and processing jobs for media-heavy workloads, including file handling, metadata indexing, and service integrations.',
    3,
    'Go,Python,PostgreSQL,Redis,MinIO'
  ),
  (
    'University Research Group',
    'Passau, Germany',
    'Mar 2021',
    'Sep 2022',
    'Student Research Assistant',
    'Implemented prototypes and tooling around data processing and service-oriented applications for academic and internal use cases.',
    4,
    'Java,Spring,Python,SQL'
  ),
  (
    'Campus IT Support',
    'Passau, Germany',
    'Oct 2019',
    'Feb 2021',
    'Technical Assistant',
    'Supported end users, automated recurring maintenance tasks, and documented operational procedures for small internal systems.',
    5,
    'Linux,Bash,Networking,Support'
  );

INSERT INTO education (institution, location, start_date, end_date, degree, thesis, description, display_order) VALUES
  (
    'University of Passau',
    'Passau, Germany',
    DATE '2021-10-01',
    DATE '2024-09-30',
    'M.Sc. Computer Science',
    'Service-oriented platform architecture for personal-scale cloud systems',
    'Graduate focus on distributed systems, backend engineering, and practical cloud operations.',
    1
  ),
  (
    'University of Passau',
    'Passau, Germany',
    DATE '2018-10-01',
    DATE '2021-09-30',
    'B.Sc. Computer Science',
    'Foundations in software engineering, algorithms, databases, and operating systems.',
    'Broad undergraduate training with emphasis on programming, systems thinking, and software design.',
    2
  ),
  (
    'Online Professional Training',
    'Remote',
    DATE '2022-01-01',
    NULL,
    'Continuous Professional Learning',
    NULL,
    'Ongoing applied learning in cloud infrastructure, DevOps workflows, frontend engineering, and AI-enabled systems.',
    3
  );

INSERT INTO projects (name, url, tech_stack, display_order) VALUES
  (
    'Personal Platform',
    'https://nzhussup.com',
    'Go,Java,React,PostgreSQL,Redis,Docker,Kubernetes',
    1
  ),
  (
    'Web Admin Panel',
    'https://admin.nzhussup.com',
    'React,TypeScript,Vite,Jest,Bootstrap',
    2
  ),
  (
    'Backend Services',
    'https://api.nzhussup.com',
    'Go,Spring Boot,PostgreSQL,Redis,OpenAPI',
    3
  ),
  (
    'Platform Infrastructure',
    'https://github.com/nzhussup/server-infra',
    'Kubernetes,k3s,NGINX,Let''s Encrypt,Bash',
    4
  ),
  (
    'CV Generator Workflow',
    'https://nzhussup.com/cv',
    'React,LLM,REST API,PDF Export',
    5
  ),
  (
    'Image Album Service',
    'https://api.nzhussup.com/v1/album',
    'Go,Redis,File Storage,REST',
    6
  );

INSERT INTO skills (category, skill_names, display_order) VALUES
  ('Backend', 'Go, Java, Spring Boot, REST, OpenAPI, JPA, JWT', 1),
  ('Frontend', 'React, TypeScript, Vite, Bootstrap, Jest, Testing Library', 2),
  ('Infrastructure', 'Docker, Kubernetes, k3s, NGINX Ingress, GitHub Actions', 3),
  ('Data', 'PostgreSQL, Redis, SQL, Data Modeling, Caching', 4),
  ('Tooling', 'Git, Bash, Linux, Maven, npm, Compose', 5),
  ('Engineering Focus', 'API Design, CI/CD, Release Management, Documentation, Maintainability', 6);

INSERT INTO certificates (name, url, display_order) VALUES
  ('AWS Cloud Practitioner', 'https://www.credly.com/', 1),
  ('Oracle Java Foundations', 'https://education.oracle.com/', 2),
  ('Docker Essentials', 'https://www.docker.com/', 3),
  ('Kubernetes Basics', 'https://kubernetes.io/', 4),
  ('GitHub Actions CI/CD', 'https://docs.github.com/actions', 5),
  ('PostgreSQL for Developers', 'https://www.postgresql.org/', 6);
