import os
import re

def fix_file(filepath):
    with open(filepath, 'r') as f:
        content = f.read()

    # Add errors import if missing and response is used
    if 'response.' in content and 'github.com/erp-cosmetics/shared/pkg/errors' not in content:
        content = re.sub(r'import \((.*?)\)', r'import (\1\n\t"github.com/erp-cosmetics/shared/pkg/errors")', content, flags=re.DOTALL)

    # Fix response.Error(c, status, msg, details)
    content = re.sub(
        r'response\.Error\(c, (http\.[A-Za-z]+|[\d]+), ([^,]+), ([^,]+|err\.Error\(\))\)',
        r'response.Error(c, errors.New("ERROR", \2, \1))',
        content
    )

    # Fix response.SuccessWithPagination(c, status, msg, data, total, page, pageSize)
    # -> SuccessWithMeta(c, data, response.NewMeta(page, pageSize, total))
    content = re.sub(
        r'response\.SuccessWithPagination\(c, [^,]+, [^,]+, ([^,]+), ([^,]+), ([^,]+), ([^)]+)\)',
        r'response.SuccessWithMeta(c, \1, response.NewMeta(\3, \4, \2))',
        content
    )

    # Fix response.Success(c, status, msg, data)
    content = re.sub(
        r'response\.Success\(c, (http\.[A-Za-z]+|[\d]+), "[^"]+", ([^)]+)\)',
        r'response.Success(c, \2)',
        content
    )
    
    # Fix response.Created(c, msg, data)
    # want (c, data)
    content = re.sub(
        r'response\.Created\(c, "[^"]+", ([^)]+)\)',
        r'response.Created(c, \1)',
        content
    )

    with open(filepath, 'w') as f:
        f.write(content)

def walk_and_fix(root_dir):
    for root, dirs, files in os.walk(root_dir):
        for file in files:
            if file.endswith('.go') and 'internal/delivery/http/handler' in root:
                fix_file(os.path.join(root, file))

if __name__ == '__main__':
    walk_and_fix('services')
