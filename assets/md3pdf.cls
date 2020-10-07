\NeedsTeXFormat{LaTeX2e}
\ProvidesClass{md3pdf}[2020/10/01 Markdown to PDF class]

\LoadClass{report}

\RequirePackage[utf8]{inputenc}
\RequirePackage{parskip}
\RequirePackage{titlesec}
\RequirePackage[table,dvipsnames]{xcolor}
\RequirePackage{listings}
\RequirePackage{float}
\RequirePackage[normalem]{ulem}

\RequirePackage{hyperref}
\hypersetup{
colorlinks=true,
linkcolor=blue,
filecolor=magenta,
urlcolor=cyan,
}

\titleformat{\paragraph}[hang]{\normalfont\normalsize\bfseries}{\theparagraph}{1em}{}
\titlespacing*{\paragraph}{0pt}{3.25ex plus 1ex minus .2ex}{1em}

\titleformat{\subparagraph}[hang]{\normalfont\normalsize\bfseries\color{gray}}{\thesubparagraph}{1em}{}
\titlespacing*{\subparagraph}{0pt}{3.25ex plus 1ex minus .2ex}{1em}

\titleformat{\chapter}[hang]{\Huge\bfseries}{}{0pt}{\Huge\bfseries}
\titleformat{\section}[hang]{\huge\bfseries}{}{0pt}{\huge\bfseries}
\titleformat{\subsection}[hang]{\Large\bfseries}{}{0pt}{\Large\bfseries}
\titleformat{\subsubsection}[hang]{\large\bfseries}{}{0pt}{\large\bfseries}

\definecolor{pblue}{rgb}{0.13,0.13,1}
\definecolor{pgreen}{rgb}{0,0.5,0}
\definecolor{pred}{rgb}{0.9,0,0}

\lstset{
    basicstyle=\ttfamily\footnotesize,
    showspaces=false,
    showtabs=false,
    breaklines=true,
    showstringspaces=false,
    breakatwhitespace=true,
    commentstyle=\color{pgreen},
    keywordstyle=\color{pblue},
    stringstyle=\color{pred},
    basicstyle=\ttfamily
}

\useunder{\uline}{\ul}{}

\setlength{\tabcolsep}{18pt}
\renewcommand{\arraystretch}{1.5}
